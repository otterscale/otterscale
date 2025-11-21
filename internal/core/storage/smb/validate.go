package smb

import (
	"crypto/tls"
	"fmt"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

const (
	ldapDefaultPort  uint16 = 389
	ldapsDefaultPort uint16 = 636
)

func (uc *UseCase) parseSearchUsername(searchUsername string) (string, error) {
	actualSearchName := searchUsername

	if strings.HasPrefix(searchUsername, "@\"") {
		// Must end with " if starts with @"
		if !strings.HasSuffix(searchUsername, "\"") {
			return "", fmt.Errorf("invalid format: missing closing quote")
		}

		// Remove @" prefix and " suffix
		actualSearchName = strings.TrimPrefix(searchUsername, "@\"")
		actualSearchName = strings.TrimSuffix(actualSearchName, "\"")

		// Reject invalid format with forward slash
		if strings.Contains(actualSearchName, "/") {
			return "", fmt.Errorf("invalid format: use backslash (\\) instead of forward slash (/)")
		}

		// If format is "DOMAIN\NAME", extract the NAME part
		if idx := strings.Index(actualSearchName, "\\"); idx != -1 {
			actualSearchName = actualSearchName[idx+1:]
		}
	}

	actualSearchName = strings.TrimSpace(actualSearchName)
	if actualSearchName == "" {
		return "", fmt.Errorf("invalid format: empty username or group name")
	}

	return actualSearchName, nil
}

func (uc *UseCase) connectLDAP(serverName string, port uint16, useTLS bool) (*ldap.Conn, error) {
	var scheme string
	var opts []ldap.DialOpt

	ldapURL := fmt.Sprintf("%s:%d", serverName, port)

	if useTLS {
		scheme = "ldaps"
		opts = append(opts, ldap.DialWithTLSConfig(&tls.Config{
			ServerName: serverName,
			MinVersion: tls.VersionTLS12,
		}))
	} else {
		scheme = "ldap"
	}

	conn, err := ldap.DialURL(fmt.Sprintf("%s://%s", scheme, ldapURL), opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to LDAP server: %w", err)
	}
	return conn, nil
}

func (uc *UseCase) ValidateSMBUser(realm, username, password, searchUsername string, useTLS bool) (EntityType, error) {
	// Parse and validate searchUsername format
	actualSearchName, err := uc.parseSearchUsername(searchUsername)
	if err != nil {
		return EntityTypeUnknown, err
	}

	port := ldapDefaultPort
	if useTLS {
		port = ldapsDefaultPort
	}

	// Connect to LDAP server
	conn, err := uc.connectLDAP(realm, port, useTLS)
	if err != nil {
		return EntityTypeUnknown, fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()

	// Format username for binding
	if !strings.Contains(username, "@") && !strings.Contains(username, "=") {
		username = fmt.Sprintf("%s@%s", username, strings.ToUpper(realm))
	}

	// Bind to LDAP
	if err = conn.Bind(username, password); err != nil {
		return EntityTypeUnknown, fmt.Errorf("authentication failed: %w", err)
	}

	// Search for user or group
	sr, err := conn.Search(ldap.NewSearchRequest(
		"DC="+strings.Join(strings.Split(strings.ToLower(realm), "."), ",DC="),
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(sAMAccountName=%s)", ldap.EscapeFilter(actualSearchName)),
		[]string{"objectClass"}, nil,
	))
	if err != nil {
		return EntityTypeUnknown, fmt.Errorf("LDAP search failed: %w", err)
	}

	if len(sr.Entries) == 0 {
		return EntityTypeUnknown, fmt.Errorf("user or group not found")
	}

	entityType := determineEntityType(sr.Entries[0].GetAttributeValues("objectClass"))
	return entityType, nil
}

func determineEntityType(objectClasses []string) EntityType {
	for _, c := range objectClasses {
		if cl := strings.ToLower(c); cl == "group" {
			return EntityTypeGroup
		}
	}
	for _, c := range objectClasses {
		if cl := strings.ToLower(c); cl == "user" || cl == "person" || cl == "inetorgperson" {
			return EntityTypeUser
		}
	}
	return EntityTypeUnknown
}
