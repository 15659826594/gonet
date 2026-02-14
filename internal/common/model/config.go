package model

// Config 配置表
type Config struct {
	Id      uint   `gorm:"column:id;primaryKey;autoIncrement"`
	Name    string `gorm:"column:name;size:30;uniqueIndex;comment:变量名"`
	Group   string `gorm:"column:group;size:30;comment:分组"`
	Title   string `gorm:"column:title;size:100;comment:变量标题"`
	Tip     string `gorm:"column:tip;size:100;comment:变量描述"`
	Type    string `gorm:"column:type;size:30;comment:类型:string,text,int,bool,array,datetime,date,file"`
	Visible string `gorm:"column:visible;size:255;comment:可见条件"`
	Value   string `gorm:"column:value;comment:变量值"`
	Content string `gorm:"column:content;comment:变量字典数据"`
	Rule    string `gorm:"column:rule;size:100;comment:验证规则"`
	Extend  string `gorm:"column:extend;size:255;comment:扩展属性"`
	Setting string `gorm:"column:setting;size:255;comment:配置"`
}
