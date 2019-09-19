package skopeo

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"net/http"

	"github.com/xelalexv/dregsy/internal/pkg/log"
)

const defaultSkopeoBinary = "skopeo"
const defaultCertsBaseDir = "/etc/skopeo/certs.d"

var skopeoBinary string
var certsBaseDir string

//
func init() {
	skopeoBinary = defaultSkopeoBinary
	certsBaseDir = defaultCertsBaseDir
}

//
type creds struct {
	Username string
	Password string
}

//
type manifest struct {
	Name     string   `json:"name"`
	RepoTags []string `json:"tags"`
}

//
func listAllTags(ref, creds, certDir string, skipTLSVerify bool) (
	[]string, error) {

	repoRef := strings.SplitN(ref, "/", 2)
	resp, err := http.Get("http://" + repoRef[0] + "/v2/" + repoRef[1] + "/tags/list")

	if err != nil  {
		return nil,
			fmt.Errorf("error listing image tags: %v", err)
	}

	return decodeManifest(resp).RepoTags, nil
}

//
func chooseOutStream(out io.Writer, verbose, isErrorStream bool) io.Writer {
	if verbose {
		if out != nil {
			return out
		}
		if log.ToTerminal {
			if isErrorStream {
				return os.Stderr
			}
			return os.Stdout
		}
	}
	return ioutil.Discard
}

//
func runSkopeo(outWr, errWr io.Writer, verbose bool, args ...string) error {

	cmd := exec.Command(skopeoBinary, args...)

	cmd.Stdout = chooseOutStream(outWr, verbose, false)
	cmd.Stderr = chooseOutStream(errWr, verbose, true)

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}

//
func decodeManifest(resp *http.Response) *manifest {
	var ret manifest

	defer closeResponse(resp)
	body, requestErr := ioutil.ReadAll(resp.Body)

	if requestErr != nil ||  json.Unmarshal(body, &ret) != nil {
		return nil
	}
	return &ret
}

func closeResponse(resp *http.Response) {
	fmt.Println("closing")
	err := resp.Body.Close()
	if err != nil {
		fmt.Errorf("error: %v\n", err)
	}
}

//
func decodeJSONAuth(authBase64 string) string {

	if authBase64 == "" {
		return ""
	}

	decoded, err := base64.StdEncoding.DecodeString(authBase64)
	if log.Error(err) {
		return ""
	}

	var ret creds
	if err := json.Unmarshal([]byte(decoded), &ret); log.Error(err) {
		return ""
	}

	return fmt.Sprintf("%s:%s", ret.Username, ret.Password)
}

//
func withoutPort(repo string) string {
	ix := strings.Index(repo, ":")
	if ix == -1 {
		return repo
	}
	return repo[:ix]
}
