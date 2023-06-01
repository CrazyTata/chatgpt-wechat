package chatAdmin

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/chatAdmin"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
    chatAdminReq "github.com/flipped-aurora/gin-vue-admin/server/model/chatAdmin/request"
)

type ApplicationConfigService struct {
}

// CreateApplicationConfig 创建ApplicationConfig记录
// Author [piexlmax](https://github.com/piexlmax)
func (applicationConfigService *ApplicationConfigService) CreateApplicationConfig(applicationConfig *chatAdmin.ApplicationConfig) (err error) {
	err = global.GVA_DB.Create(applicationConfig).Error
	return err
}

// DeleteApplicationConfig 删除ApplicationConfig记录
// Author [piexlmax](https://github.com/piexlmax)
func (applicationConfigService *ApplicationConfigService)DeleteApplicationConfig(applicationConfig chatAdmin.ApplicationConfig) (err error) {
	err = global.GVA_DB.Delete(&applicationConfig).Error
	return err
}

// DeleteApplicationConfigByIds 批量删除ApplicationConfig记录
// Author [piexlmax](https://github.com/piexlmax)
func (applicationConfigService *ApplicationConfigService)DeleteApplicationConfigByIds(ids request.IdsReq) (err error) {
	err = global.GVA_DB.Delete(&[]chatAdmin.ApplicationConfig{},"id in ?",ids.Ids).Error
	return err
}

// UpdateApplicationConfig 更新ApplicationConfig记录
// Author [piexlmax](https://github.com/piexlmax)
func (applicationConfigService *ApplicationConfigService)UpdateApplicationConfig(applicationConfig chatAdmin.ApplicationConfig) (err error) {
	err = global.GVA_DB.Save(&applicationConfig).Error
	return err
}

// GetApplicationConfig 根据id获取ApplicationConfig记录
// Author [piexlmax](https://github.com/piexlmax)
func (applicationConfigService *ApplicationConfigService)GetApplicationConfig(id uint) (applicationConfig chatAdmin.ApplicationConfig, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&applicationConfig).Error
	return
}

// GetApplicationConfigInfoList 分页获取ApplicationConfig记录
// Author [piexlmax](https://github.com/piexlmax)
func (applicationConfigService *ApplicationConfigService)GetApplicationConfigInfoList(info chatAdminReq.ApplicationConfigSearch) (list []chatAdmin.ApplicationConfig, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&chatAdmin.ApplicationConfig{})
    var applicationConfigs []chatAdmin.ApplicationConfig
    // 如果有条件搜索 下方会自动创建搜索语句
    if info.StartCreatedAt !=nil && info.EndCreatedAt !=nil {
     db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
    }
    if info.AgentName != "" {
        db = db.Where("agent_name = ?",info.AgentName)
    }
    if info.Model != "" {
        db = db.Where("model = ?",info.Model)
    }
	err = db.Count(&total).Error
	if err!=nil {
    	return
    }

	err = db.Limit(limit).Offset(offset).Find(&applicationConfigs).Error
	return  applicationConfigs, total, err
}
