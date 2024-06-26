// Code generated by MockGen. DO NOT EDIT.
// Source: repository/repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	sql "database/sql"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	models "github.com/mashmorsik/banners-service/pkg/models"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// AddNewFeature mocks base method.
func (m *MockRepository) AddNewFeature(banner *models.Banner) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddNewFeature", banner)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddNewFeature indicates an expected call of AddNewFeature.
func (mr *MockRepositoryMockRecorder) AddNewFeature(banner interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddNewFeature", reflect.TypeOf((*MockRepository)(nil).AddNewFeature), banner)
}

// AddNewTag mocks base method.
func (m *MockRepository) AddNewTag(banner *models.Banner) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddNewTag", banner)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddNewTag indicates an expected call of AddNewTag.
func (mr *MockRepositoryMockRecorder) AddNewTag(banner interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddNewTag", reflect.TypeOf((*MockRepository)(nil).AddNewTag), banner)
}

// CheckTagFeatureOverlap mocks base method.
func (m *MockRepository) CheckTagFeatureOverlap(b *models.Banner) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckTagFeatureOverlap", b)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckTagFeatureOverlap indicates an expected call of CheckTagFeatureOverlap.
func (mr *MockRepositoryMockRecorder) CheckTagFeatureOverlap(b interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckTagFeatureOverlap", reflect.TypeOf((*MockRepository)(nil).CheckTagFeatureOverlap), b)
}

// Create mocks base method.
func (m *MockRepository) Create(b *models.Banner) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", b)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockRepositoryMockRecorder) Create(b interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRepository)(nil).Create), b)
}

// CreateBanner mocks base method.
func (m *MockRepository) CreateBanner(tx *sql.Tx, b *models.Banner) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateBanner", tx, b)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateBanner indicates an expected call of CreateBanner.
func (mr *MockRepositoryMockRecorder) CreateBanner(tx, b interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBanner", reflect.TypeOf((*MockRepository)(nil).CreateBanner), tx, b)
}

// CreateContent mocks base method.
func (m *MockRepository) CreateContent(tx *sql.Tx, b *models.Banner) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateContent", tx, b)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateContent indicates an expected call of CreateContent.
func (mr *MockRepositoryMockRecorder) CreateContent(tx, b interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateContent", reflect.TypeOf((*MockRepository)(nil).CreateContent), tx, b)
}

// CreateFeatureTags mocks base method.
func (m *MockRepository) CreateFeatureTags(tx *sql.Tx, b *models.Banner) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateFeatureTags", tx, b)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateFeatureTags indicates an expected call of CreateFeatureTags.
func (mr *MockRepositoryMockRecorder) CreateFeatureTags(tx, b interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateFeatureTags", reflect.TypeOf((*MockRepository)(nil).CreateFeatureTags), tx, b)
}

// Delete mocks base method.
func (m *MockRepository) Delete(bannerID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", bannerID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockRepositoryMockRecorder) Delete(bannerID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRepository)(nil).Delete), bannerID)
}

// GetBannerActiveVersions mocks base method.
func (m *MockRepository) GetBannerActiveVersions() (map[int]int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBannerActiveVersions")
	ret0, _ := ret[0].(map[int]int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBannerActiveVersions indicates an expected call of GetBannerActiveVersions.
func (mr *MockRepositoryMockRecorder) GetBannerActiveVersions() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBannerActiveVersions", reflect.TypeOf((*MockRepository)(nil).GetBannerActiveVersions))
}

// GetForAdmin mocks base method.
func (m *MockRepository) GetForAdmin(b *models.Banner, limit, offset int) ([]*models.Banner, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetForAdmin", b, limit, offset)
	ret0, _ := ret[0].([]*models.Banner)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetForAdmin indicates an expected call of GetForAdmin.
func (mr *MockRepositoryMockRecorder) GetForAdmin(b, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetForAdmin", reflect.TypeOf((*MockRepository)(nil).GetForAdmin), b, limit, offset)
}

// GetForUser mocks base method.
func (m *MockRepository) GetForUser(b *models.Banner) (*models.Banner, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetForUser", b)
	ret0, _ := ret[0].(*models.Banner)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetForUser indicates an expected call of GetForUser.
func (mr *MockRepositoryMockRecorder) GetForUser(b interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetForUser", reflect.TypeOf((*MockRepository)(nil).GetForUser), b)
}

// MergeUpdateVersion mocks base method.
func (m *MockRepository) MergeUpdateVersion(tx *sql.Tx, b *models.Banner, lastVersion int) (*models.Banner, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MergeUpdateVersion", tx, b, lastVersion)
	ret0, _ := ret[0].(*models.Banner)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MergeUpdateVersion indicates an expected call of MergeUpdateVersion.
func (mr *MockRepositoryMockRecorder) MergeUpdateVersion(tx, b, lastVersion interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MergeUpdateVersion", reflect.TypeOf((*MockRepository)(nil).MergeUpdateVersion), tx, b, lastVersion)
}

// SetVersionActive mocks base method.
func (m *MockRepository) SetVersionActive(bannerID, version int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetVersionActive", bannerID, version)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetVersionActive indicates an expected call of SetVersionActive.
func (mr *MockRepositoryMockRecorder) SetVersionActive(bannerID, version interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetVersionActive", reflect.TypeOf((*MockRepository)(nil).SetVersionActive), bannerID, version)
}

// Update mocks base method.
func (m *MockRepository) Update(b *models.Banner) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", b)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockRepositoryMockRecorder) Update(b interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockRepository)(nil).Update), b)
}

// UpdateBanner mocks base method.
func (m *MockRepository) UpdateBanner(tx *sql.Tx, b *models.Banner, lastVersion int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBanner", tx, b, lastVersion)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateBanner indicates an expected call of UpdateBanner.
func (mr *MockRepositoryMockRecorder) UpdateBanner(tx, b, lastVersion interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBanner", reflect.TypeOf((*MockRepository)(nil).UpdateBanner), tx, b, lastVersion)
}

// UpdateBannerContent mocks base method.
func (m *MockRepository) UpdateBannerContent(tx *sql.Tx, b *models.Banner) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBannerContent", tx, b)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateBannerContent indicates an expected call of UpdateBannerContent.
func (mr *MockRepositoryMockRecorder) UpdateBannerContent(tx, b interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBannerContent", reflect.TypeOf((*MockRepository)(nil).UpdateBannerContent), tx, b)
}

// UpdateFeatureTag mocks base method.
func (m *MockRepository) UpdateFeatureTag(tx *sql.Tx, b *models.Banner) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateFeatureTag", tx, b)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateFeatureTag indicates an expected call of UpdateFeatureTag.
func (mr *MockRepositoryMockRecorder) UpdateFeatureTag(tx, b interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateFeatureTag", reflect.TypeOf((*MockRepository)(nil).UpdateFeatureTag), tx, b)
}
