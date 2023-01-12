package actions

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/debdutdeb/rocketchatctl/backend"
	"github.com/debdutdeb/rocketchatctl/deployment"
	"golang.org/x/mod/semver"
)

func InstallAllResources(ctx context.Context, backend backend.Backend, opts deployment.HostDeploymentOptions) error {
	if err := verifyRocketChatRelease(opts.Version); err != nil {
		return err
	}
	// if err := backend.PullArtifacts(ctx, opts); err != nil {
	// 	return err
	// }
	if err := backend.StartDeployment(ctx, opts); err != nil {
		return err
	}
	// if err := backend.ConfigureHost(ctx, opts); err != nil {
	// 	return err
	// }
	return nil
}

type releaseInfo struct {
	Tag string `json:"tag,omitempty"`
}

func verifyRocketChatRelease(version string) error {
	if !semver.IsValid("v" + version) {
		return fmt.Errorf("%s is not a valid Rocket.Chat release", version)
	}
	resp, err := http.Get(fmt.Sprintf("https://releases.rocket.chat/%s/info", version))
	if err != nil {
		return fmt.Errorf("failed to verify given release %v: %v", version, err)
	}
	defer resp.Body.Close()
	info := releaseInfo{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to verify given release %v: %v", version, err)
	}
	if err = json.Unmarshal(body, &info); err != nil {
		return fmt.Errorf("failed to verify given release %v: %v", version, err)
	}
	log.Printf("info: %v\n", info)
	if info.Tag == "" {
		return fmt.Errorf("failed to get release information for version string %v", version)
	}
	return nil
}
