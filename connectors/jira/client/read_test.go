package client

import (
	"context"
	"fmt"
	"testing"

	"github.com/openhdc/openhdc"
	pb "github.com/openhdc/openhdc/api/connector/v1"
)

func TestClient_readProjects(t *testing.T) {
	type args struct {
		msgs chan<- *pb.Message
		rdr  *openhdc.Reader
	}
	tests := []struct {
		name    string
		c       *Client
		args    args
		wantErr bool
	}{
		// {
		// 	name: "nil client",
		// 	c:    nil,
		// 	args: args{
		// 		msgs: make(chan *pb.Message, 1),
		// 		rdr:  &openhdc.Reader{},
		// 	},
		// 	wantErr: true,
		// },
		// {
		// 	name: "nil reader",
		// 	c:    &Client{},
		// 	args: args{
		// 		msgs: make(chan *pb.Message, 1),
		// 		rdr:  nil,
		// 	},
		// 	wantErr: true,
		// },
		// {
		// 	name: "empty messages channel",
		// 	c:    &Client{},
		// 	args: args{
		// 		msgs: nil,
		// 		rdr:  &openhdc.Reader{},
		// 	},
		// 	wantErr: true,
		// },
		// {
		// 	name: "valid client and reader",
		// 	c:    &Client{},
		// 	args: args{
		// 		msgs: make(chan *pb.Message, 1),
		// 		rdr:  &openhdc.Reader{},
		// 	},
		// 	wantErr: false,
		// },
		// {
		// 	name: "closed messages channel",
		// 	c:    &Client{},
		// 	args: func() args {
		// 		ch := make(chan *pb.Message, 1)
		// 		close(ch)
		// 		return args{
		// 			msgs: ch,
		// 			rdr:  &openhdc.Reader{},
		// 		}
		// 	}(),
		// 	wantErr: true,
		// },
		// {
		// 	name: "reader with nil config",
		// 	c:    &Client{},
		// 	args: args{
		// 		msgs: make(chan *pb.Message, 1),
		// 		rdr:  nil,
		// 	},
		// 	wantErr: true,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Printf("tt.args.msgs:%v, tt.args.rdr:%v, tt.wantErr:%v,", tt.args.msgs, tt.args.rdr, tt.wantErr)
			if err := tt.c.readProjects(tt.args.msgs, tt.args.rdr); (err != nil) != tt.wantErr {
				t.Errorf("Client.readProjects() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_readIssueFields(t *testing.T) {
	type args struct {
		msgs chan<- *pb.Message
		rdr  *openhdc.Reader
	}
	tests := []struct {
		name    string
		c       *Client
		args    args
		wantErr bool
	}{
		// {
		// 	name: "nil client",
		// 	c:    nil,
		// 	args: args{
		// 		msgs: make(chan *pb.Message, 1),
		// 		rdr:  &openhdc.Reader{},
		// 	},
		// 	wantErr: true,
		// },
		// {
		// 	name: "nil reader",
		// 	c:    &Client{},
		// 	args: args{
		// 		msgs: make(chan *pb.Message, 1),
		// 		rdr:  nil,
		// 	},
		// 	wantErr: true,
		// },
		// {
		// 	name: "empty messages channel",
		// 	c:    &Client{},
		// 	args: args{
		// 		msgs: nil,
		// 		rdr:  &openhdc.Reader{},
		// 	},
		// 	wantErr: true,
		// },
		// {
		// 	name: "valid client and reader",
		// 	c:    &Client{},
		// 	args: args{
		// 		msgs: make(chan *pb.Message, 1),
		// 		rdr:  &openhdc.Reader{},
		// 	},
		// 	wantErr: false,
		// },
		// {
		// 	name: "closed messages channel",
		// 	c:    &Client{},
		// 	args: func() args {
		// 		ch := make(chan *pb.Message, 1)
		// 		close(ch)
		// 		return args{
		// 			msgs: ch,
		// 			rdr:  &openhdc.Reader{},
		// 		}
		// 	}(),
		// 	wantErr: true,
		// },
		// {
		// 	name: "reader with nil config",
		// 	c:    &Client{},
		// 	args: args{
		// 		msgs: make(chan *pb.Message, 1),
		// 		rdr:  nil,
		// 	},
		// 	wantErr: true,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.readIssueFields(tt.args.msgs, tt.args.rdr); (err != nil) != tt.wantErr {
				t.Errorf("Client.readIssueFields() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_readIssues(t *testing.T) {
	type args struct {
		pj   string
		sd   string
		msgs chan<- *pb.Message
		rdr  *openhdc.Reader
	}
	tests := []struct {
		name    string
		c       *Client
		args    args
		wantErr bool
	}{
		// {
		// 	name: "nil client",
		// 	c:    nil,
		// 	args: args{
		// 		pj:   "PROJ",
		// 		sd:   "2024-01-01",
		// 		msgs: make(chan *pb.Message, 1),
		// 		rdr:  &openhdc.Reader{},
		// 	},
		// 	wantErr: true,
		// },
		// {
		// 	name: "nil reader",
		// 	c:    &Client{},
		// 	args: args{
		// 		pj:   "PROJ",
		// 		sd:   "2024-01-01",
		// 		msgs: make(chan *pb.Message, 1),
		// 		rdr:  nil,
		// 	},
		// 	wantErr: true,
		// },
		// {
		// 	name: "empty messages channel",
		// 	c:    &Client{},
		// 	args: args{
		// 		pj:   "PROJ",
		// 		sd:   "2024-01-01",
		// 		msgs: nil,
		// 		rdr:  &openhdc.Reader{},
		// 	},
		// 	wantErr: true,
		// },
		// {
		// 	name: "valid client and reader",
		// 	c:    &Client{},
		// 	args: args{
		// 		pj:   "PROJ",
		// 		sd:   "2024-01-01",
		// 		msgs: make(chan *pb.Message, 1),
		// 		rdr:  &openhdc.Reader{},
		// 	},
		// 	wantErr: false,
		// },
		// {
		// 	name: "closed messages channel",
		// 	c:    &Client{},
		// 	args: func() args {
		// 		ch := make(chan *pb.Message, 1)
		// 		close(ch)
		// 		return args{
		// 			pj:   "PROJ",
		// 			sd:   "2024-01-01",
		// 			msgs: ch,
		// 			rdr:  &openhdc.Reader{},
		// 		}
		// 	}(),
		// 	wantErr: true,
		// },
		// {
		// 	name: "empty project key",
		// 	c:    &Client{},
		// 	args: args{
		// 		pj:   "",
		// 		sd:   "2024-01-01",
		// 		msgs: make(chan *pb.Message, 1),
		// 		rdr:  &openhdc.Reader{},
		// 	},
		// 	wantErr: true,
		// },
		// {
		// 	name: "empty start date",
		// 	c:    &Client{},
		// 	args: args{
		// 		pj:   "PROJ",
		// 		sd:   "",
		// 		msgs: make(chan *pb.Message, 1),
		// 		rdr:  &openhdc.Reader{},
		// 	},
		// 	wantErr: true,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.readIssues(tt.args.pj, tt.args.sd, tt.args.msgs, tt.args.rdr); (err != nil) != tt.wantErr {
				t.Errorf("Client.readIssues() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_Read(t *testing.T) {
	type args struct {
		ctx  context.Context
		msgs chan<- *pb.Message
		rdr  *openhdc.Reader
	}
	tests := []struct {
		name    string
		c       *Client
		args    args
		wantErr bool
	}{
		// {
		// 	name: "nil client",
		// 	c:    nil,
		// 	args: args{
		// 		ctx:  context.Background(),
		// 		msgs: make(chan *pb.Message, 1),
		// 		rdr:  &openhdc.Reader{},
		// 	},
		// 	wantErr: true,
		// },
		// {
		// 	name: "nil reader",
		// 	c:    &Client{},
		// 	args: args{
		// 		ctx:  context.Background(),
		// 		msgs: make(chan *pb.Message, 1),
		// 		rdr:  nil,
		// 	},
		// 	wantErr: true,
		// },
		// {
		// 	name: "nil messages channel",
		// 	c:    &Client{},
		// 	args: args{
		// 		ctx:  context.Background(),
		// 		msgs: nil,
		// 		rdr:  &openhdc.Reader{},
		// 	},
		// 	wantErr: true,
		// },
		// {
		// 	name: "closed messages channel",
		// 	c:    &Client{},
		// 	args: func() args {
		// 		ch := make(chan *pb.Message, 1)
		// 		close(ch)
		// 		return args{
		// 			ctx:  context.Background(),
		// 			msgs: ch,
		// 			rdr:  &openhdc.Reader{},
		// 		}
		// 	}(),
		// 	wantErr: true,
		// },
		// {
		// 	name: "empty namespace in options",
		// 	c: &Client{
		// 		opts: options{server: ""},
		// 	},
		// 	args: args{
		// 		ctx:  context.Background(),
		// 		msgs: make(chan *pb.Message, 1),
		// 		rdr:  &openhdc.Reader{},
		// 	},
		// 	wantErr: true,
		// },
		// You can add more cases for error returned by pool.OpenWithContext, pool.BeginTx, newTables, etc.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.Read(tt.args.ctx, tt.args.msgs, tt.args.rdr); (err != nil) != tt.wantErr {
				t.Errorf("Client.Read() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}


