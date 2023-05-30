package chatAdmin

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/chatAdmin"
	chatAdminReq "github.com/flipped-aurora/gin-vue-admin/server/model/chatAdmin/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type CustomerConfigService struct {
}

// CreateCustomerConfig 创建CustomerConfig记录
// Author [piexlmax](https://github.com/piexlmax)
func (customerConfigService *CustomerConfigService) CreateCustomerConfig(customerConfig *chatAdmin.CustomerConfig) (err error) {
	err = global.GVA_DB.Create(customerConfig).Error
	return err
}

// DeleteCustomerConfig 删除CustomerConfig记录
// Author [piexlmax](https://github.com/piexlmax)
func (customerConfigService *CustomerConfigService) DeleteCustomerConfig(customerConfig chatAdmin.CustomerConfig) (err error) {
	err = global.GVA_DB.Delete(&customerConfig).Error
	return err
}

// DeleteCustomerConfigByIds 批量删除CustomerConfig记录
// Author [piexlmax](https://github.com/piexlmax)
func (customerConfigService *CustomerConfigService) DeleteCustomerConfigByIds(ids request.IdsReq) (err error) {
	err = global.GVA_DB.Delete(&[]chatAdmin.CustomerConfig{}, "id in ?", ids.Ids).Error
	return err
}

// UpdateCustomerConfig 更新CustomerConfig记录
// Author [piexlmax](https://github.com/piexlmax)
func (customerConfigService *CustomerConfigService) UpdateCustomerConfig(customerConfig chatAdmin.CustomerConfig) (err error) {
	err = global.GVA_DB.Save(&customerConfig).Error
	return err
}

// GetCustomerConfig 根据id获取CustomerConfig记录
// Author [piexlmax](https://github.com/piexlmax)
func (customerConfigService *CustomerConfigService) GetCustomerConfig(id uint) (customerConfig chatAdmin.CustomerConfig, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&customerConfig).Error
	return
}

// GetCustomerConfigInfoList 分页获取CustomerConfig记录
// Author [piexlmax](https://github.com/piexlmax)
func (customerConfigService *CustomerConfigService) GetCustomerConfigInfoList(info chatAdminReq.CustomerConfigSearch) (list []chatAdmin.CustomerConfig, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&chatAdmin.CustomerConfig{})
	var customerConfigs []chatAdmin.CustomerConfig
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.KfName != "" {
		db = db.Where("kf_name = ?", info.KfName)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Find(&customerConfigs).Error
	return customerConfigs, total, err
}
