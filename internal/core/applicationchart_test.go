// applicationchart_test.go
package core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/repo"
)

type mockChartRepo struct {
	mock.Mock
}

func (m *mockChartRepo) List(ctx context.Context) ([]Chart, error) {
	args := m.Called(ctx)
	return args.Get(0).([]Chart), args.Error(1)
}

func (m *mockChartRepo) Show(chartRef string, format action.ShowOutputFormat) (string, error) {
	args := m.Called(chartRef, format)
	return args.String(0), args.Error(1)
}

// Add other mocks...

func TestApplicationUseCase_ListCharts(t *testing.T) {
	ctx := context.Background()
	uc := &ApplicationUseCase{
		chart: &mockChartRepo{},
	}

	mockChart := uc.chart.(*mockChartRepo)
	expectedCharts := []Chart{{Name: "test"}}
	mockChart.On("List", ctx).Return(expectedCharts, nil)

	charts, err := uc.ListCharts(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expectedCharts, charts)

	mockChart.AssertExpectations(t)
}

func TestApplicationUseCase_GetChart(t *testing.T) {
	ctx := context.Background()
	uc := &ApplicationUseCase{
		chart: &mockChartRepo{},
	}

	mockChart := uc.chart.(*mockChartRepo)
	expectedChart := &Chart{Name: "test", Versions: repo.ChartVersions{}}
	mockChart.On("List", ctx).Return([]Chart{{Name: "test", Versions: repo.ChartVersions{}}}, nil)

	chart, err := uc.GetChart(ctx, "test")
	assert.NoError(t, err)
	assert.Equal(t, expectedChart, chart)

	mockChart.AssertExpectations(t)
}

// Add tests for GetChartMetadataFromApplication, GetChartMetadata...
