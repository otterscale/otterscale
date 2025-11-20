package smb

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

const ldapDefaultPort uint16 = 389

const (
	EntityTypeUnknown int = iota
	EntityTypeUser
	EntityTypeGroup
)

type ADValidateResult struct {
	Valid      bool
	EntityType int
	Message    string
}

func (uc *UseCase) parseSearchUsername(searchUsername string) (string, *ADValidateResult) {
	result := &ADValidateResult{EntityType: EntityTypeUnknown}
	actualSearchName := searchUsername

	if strings.HasPrefix(searchUsername, "@\"") {
		// Must end with " if starts with @"
		if !strings.HasSuffix(searchUsername, "\"") {
			result.Valid = false
			result.Message = "invalid format: missing closing quote"
			return "", result
		}

		// Remove @" prefix and " suffix
		actualSearchName = strings.TrimPrefix(searchUsername, "@\"")
		actualSearchName = strings.TrimSuffix(actualSearchName, "\"")

		// Reject invalid format with forward slash
		if strings.Contains(actualSearchName, "/") {
			result.Valid = false
			result.Message = "invalid format: use backslash (\\) instead of forward slash (/)"
			return "", result
		}

		// If format is "DOMAIN\NAME", extract the NAME part
		if idx := strings.Index(actualSearchName, "\\"); idx != -1 {
			actualSearchName = actualSearchName[idx+1:]
		}
	}

	actualSearchName = strings.TrimSpace(actualSearchName)
	if actualSearchName == "" {
		result.Valid = false
		result.Message = "invalid format: empty username or group name"
		return "", result
	}

	return actualSearchName, nil
}

func (uc *UseCase) resolveLDAPServer(ctx context.Context, realm string) (serverName string, port uint16, err error) {
	resolver := &net.Resolver{}
	_, srvs, err := resolver.LookupSRV(ctx, "ldap", "tcp", realm)
	if err == nil && len(srvs) > 0 {
		return strings.TrimSuffix(srvs[0].Target, "."), srvs[0].Port, nil
	}

	if addrs, err := resolver.LookupHost(ctx, realm); err == nil && len(addrs) > 0 {
		return realm, ldapDefaultPort, nil
	}

	return "", 0, fmt.Errorf("failed to lookup LDAP server: unable to resolve %s", realm)
}

func (uc *UseCase) connectLDAP(serverName string, port uint16, useTLS bool) (*ldap.Conn, error) {
	ldapURL := fmt.Sprintf("%s:%d", serverName, port)

	if useTLS {
		conn, err := ldap.DialURL(fmt.Sprintf("ldaps://%s", ldapURL),
			ldap.DialWithTLSConfig(&tls.Config{
				ServerName: serverName,
				MinVersion: tls.VersionTLS12,
			}))
		if err != nil {
			return nil, fmt.Errorf("failed to connect to LDAP server: %w", err)
		}
		return conn, nil
	}

	conn, err := ldap.DialURL(fmt.Sprintf("ldap://%s", ldapURL))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to LDAP server: %w", err)
	}
	return conn, nil
}

func (uc *UseCase) ValidateSMBUser(ctx context.Context, realm, username, password, searchUsername string, useTLS bool) (*ADValidateResult, error) {
	result := &ADValidateResult{EntityType: EntityTypeUnknown}
	searchUsername = strings.TrimSpace(searchUsername)

	// Parse and validate searchUsername format
	actualSearchName, parseResult := uc.parseSearchUsername(searchUsername)
	if parseResult != nil {
		return parseResult, nil
	}

	// Resolve LDAP server
	serverName, port, err := uc.resolveLDAPServer(ctx, realm)
	if err != nil {
		result.Message = "LDAP server not found"
		return result, err
	}

	// Connect to LDAP server
	conn, err := uc.connectLDAP(serverName, port, useTLS)
	if err != nil {
		result.Message = "failed to connect"
		return result, err
	}
	defer conn.Close()

	// Format username for binding
	if !strings.Contains(username, "@") && !strings.Contains(username, "=") {
		username = fmt.Sprintf("%s@%s", username, strings.ToUpper(realm))
	}

	// Bind to LDAP
	if err = conn.Bind(username, password); err != nil {
		result.Message = "authentication failed"
		return result, fmt.Errorf("LDAP bind failed")
	}

	// Search for user or group
	sr, err := conn.Search(ldap.NewSearchRequest(
		"DC="+strings.Join(strings.Split(strings.ToLower(realm), "."), ",DC="),
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(sAMAccountName=%s)", ldap.EscapeFilter(actualSearchName)),
		[]string{"objectClass"}, nil,
	))
	if err != nil {
		result.Message = "search failed"
		return result, fmt.Errorf("LDAP search failed")
	}

	if len(sr.Entries) == 0 {
		result.Message = "user or group not found"
		return result, nil
	}

	result.Valid = true
	result.EntityType = determineEntityType(sr.Entries[0].GetAttributeValues("objectClass"))
	result.Message = "validation successful"

	return result, nil
}

func determineEntityType(objectClasses []string) int {
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
