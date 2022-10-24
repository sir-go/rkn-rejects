package fw

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"golang.org/x/net/context"
)

const (
	TimeoutExec = 2
	//MaxElemensAmount = 100
)

// NfSetElement key - ip or subnet, val - timeout
type (
	NfTables struct {
		TableName string
	}

	NfSetElement struct {
		Ip      string
		Timeout uint32
		Comment string
	}
)

func (el *NfSetElement) String() string {
	r := el.Ip
	if el.Timeout > 1 {
		r += fmt.Sprintf(" timeout %ds", el.Timeout)
	}
	if el.Comment == "" {
		return r
	}
	if len(el.Comment) > 20 {
		el.Comment = el.Comment[:20]
	}
	return fmt.Sprintf("%s comment \"%s\"", r, el.Comment)
}

func makeCtx(sec int) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(),
		time.Second*time.Duration(sec))
}

func execNft(cmd, tableName, setName, elBlock string) error {
	ctx, cancel := makeCtx(TimeoutExec)
	defer cancel()
	out, err := exec.CommandContext(ctx, "nft",
		cmd, "element", tableName, setName, elBlock,
	).CombinedOutput()
	if err != nil {
		if strings.Contains(string(out), "not found in set") {
			return nil
		}
		return fmt.Errorf("%s elements from %s: out: %s, err: %v",
			cmd, setName, out, err)
	}
	return nil
}

func Del(tableName, setName string, el NfSetElement) error {
	err := execNft("delete",
		tableName, setName, "{ "+el.Ip+" }")
	if err != nil {
		if !strings.Contains(err.Error(),
			"No such file or directory") {
			return err
		}
	}
	return nil
}

func Add(tableName, setName string, el NfSetElement) error {
	return execNft("add",
		tableName, setName, "{ "+el.String()+" }")
}

func FlushSet(tableName, setName string) error {
	ctx, cancel := makeCtx(TimeoutExec)
	defer cancel()
	out, err := exec.CommandContext(ctx, "nft",
		"flush", "set", tableName, setName).CombinedOutput()
	if err != nil {
		return fmt.Errorf("nft flush set %s %s, out: %s, err: %v",
			tableName, setName, out, err)
	}
	return nil
}

func Apply(fileName string) error {
	ctx, cancel := makeCtx(TimeoutExec)
	defer cancel()
	out, err := exec.CommandContext(ctx, "nft",
		"-f", fileName).CombinedOutput()
	if err != nil {
		return fmt.Errorf("nft -f %s, out: %s, err: %v",
			fileName, out, err)
	}
	return nil
}
