package wikiclient

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	pb "github.com/omustardo/wikitree-api-client/go/proto"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

// defaultEndpoint is the WikiTree API.
const defaultEndpoint = "https://api.wikitree.com/api.php"

// New returns a client capable of interacting with WikiTree.
func New() (pb.WikiTreeClient, error) {
	return &wikiClientImpl{
		endpoint: defaultEndpoint,
	}, nil
}

type wikiClientImpl struct {
	// URL to query. It is a field rather than a hard-coded value
	// so that tests can override it to verify that HTTP queries
	// are sent as expected.
	endpoint string
}

// GetProfile returns information about a profile.
// https://github.com/wikitree/wikitree-api/blob/main/getProfile.md
func (c *wikiClientImpl) GetProfile(_ context.Context, req *pb.GetProfileRequest, _ ...grpc.CallOption) (*pb.GetProfileResponse, error) {
	jsonResp, err := c.executeQuery(genGetProfileParams(req))
	if err != nil {
		return nil, errors.Wrap(err, "error in GetProfile")
	}
	fmt.Println("== Raw JSON Response")
	fmt.Println(string(jsonResp))

	resp := new(pb.GetProfileResponse)
	decoder := protojson.UnmarshalOptions{
		DiscardUnknown: true,
		AllowPartial:   true,
	}
	if err := decoder.Unmarshal(jsonResp, resp); err != nil {
		return nil, errors.Wrap(err, "error unmarshalling GetProfile response")
	}
	return resp, nil
}

// executeQuery does a GET request to WikiTree with the provided parameter string.
// Input is expected to be something like: "?action=getProfile&key=Mustardo-1"
// Output is the JSON returned by the API, suitable for decoding into a proto.
func (c *wikiClientImpl) executeQuery(params string) ([]byte, error) {
	httpResp, err := http.Get(c.endpoint + params)
	if err != nil {
		return nil, errors.Wrap(err, "error calling HTTP endpoint")
	}
	jsonResp, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error getting JSON response")
	}
	// The JSON response has square brackets at the start and end. This makes it a JSON array.
	// Unfortunately proto conversion cannot work with a top level array, as in protobufs
	// everything is nested within a struct.
	// As a best-effort approach, we strip off the surrounding brackets.
	if len(jsonResp) >= 2 {
		jsonResp = jsonResp[1 : len(jsonResp)-1]
	}
	return jsonResp, nil
}

// genGetProfileParams concatenates fields in the input message into URL parameters,
// suitable for concatenation with the base WikiTree URL.
// Sample return: "?action=getProfile&key=Mustardo-1"
func genGetProfileParams(req *pb.GetProfileRequest) string {
	var buf bytes.Buffer
	buf.WriteString("?action=")
	buf.WriteString(req.GetAction())
	if req.GetKey() != "" {
		buf.WriteString("&key=")
		buf.WriteString(req.GetKey())
	}
	if req.GetFields() != "" {
		buf.WriteString("&fields=")
		buf.WriteString(req.GetFields())
	}
	if req.GetBioFormat() != "" {
		buf.WriteString("&bioFormat=")
		buf.WriteString(req.GetKey())
	}
	if req.GetResolveRedirect() != "" && req.GetResolveRedirect() != "1" {
		buf.WriteString("&resolveRedirect=")
		buf.WriteString(req.GetResolveRedirect())
	}
	return buf.String()
}
