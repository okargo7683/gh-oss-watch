package cmd_test

import (
	"fmt"
	"testing"

	"github.com/jackchuka/gh-oss-watch/cmd"
	"github.com/jackchuka/gh-oss-watch/services"
	mock_services "github.com/jackchuka/gh-oss-watch/services/mock"
	"go.uber.org/mock/gomock"
)

func TestHandleInit_Success(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockConfig := mock_services.NewMockConfigService(ctrl)
	mockOutput := mock_services.NewMockOutput(ctrl)

	// Set up expectations
	mockConfig.EXPECT().Load().Return(&services.Config{Repos: []services.RepoConfig{}}, nil)
	mockConfig.EXPECT().GetConfigPath().Return("/mock/config.yaml", nil)
	mockConfig.EXPECT().Save(gomock.Any()).Return(nil)
	mockOutput.EXPECT().Printf(gomock.Any(), gomock.Any()).AnyTimes()
	mockOutput.EXPECT().Println(gomock.Any()).AnyTimes()

	err := cmd.HandleInit(mockConfig, mockOutput)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestHandleInit_LoadError(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockConfig := mock_services.NewMockConfigService(ctrl)
	mockOutput := mock_services.NewMockOutput(ctrl)

	// Set up expectation for Load to return error
	mockConfig.EXPECT().Load().Return(nil, fmt.Errorf("load failed"))

	err := cmd.HandleInit(mockConfig, mockOutput)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "load failed" {
		t.Errorf("Expected 'load failed', got %v", err)
	}
}

func TestHandleInit_SaveError(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockConfig := mock_services.NewMockConfigService(ctrl)
	mockOutput := mock_services.NewMockOutput(ctrl)

	// Set up expectations
	mockConfig.EXPECT().Load().Return(&services.Config{Repos: []services.RepoConfig{}}, nil)
	mockConfig.EXPECT().GetConfigPath().Return("/mock/config.yaml", nil)
	mockConfig.EXPECT().Save(gomock.Any()).Return(fmt.Errorf("save failed"))

	err := cmd.HandleInit(mockConfig, mockOutput)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "save failed" {
		t.Errorf("Expected 'save failed', got %v", err)
	}
}
