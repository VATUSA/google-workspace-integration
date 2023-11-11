package workspace_helper

import (
	"context"
	"github.com/VATUSA/google-workspace-integration/internal/config"
	"golang.org/x/oauth2/google"
	dir "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/option"
	"net/http"
)

func getClient() *http.Client {
	b := []byte(config.SERVICE_ACCOUNT_JSON)
	//b, err := os.ReadFile("E:\\Creds\\VATUSA\\vatusa-integration-e6cc974dab1c.json")
	//if err != nil {
	//	print(err.Error())
	//}
	conf, _ := google.JWTConfigFromJSON(b, dir.AdminDirectoryGroupScope, dir.AdminDirectoryUserScope)
	client := conf.Client(context.Background())
	return client
}

func GetService() (*dir.Service, error) {
	ctx := context.Background()
	client := getClient()
	srv, err := dir.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}
	return srv, nil
}
