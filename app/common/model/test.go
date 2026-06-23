package model

// Test 测试表
type Test struct {
	Id           uint    `gorm:"column:id;primaryKey;autoIncrement;comment:ID"`
	UserId       int     `gorm:"column:user_id;default:0;comment:会员ID"`
	AdminId      int     `gorm:"column:admin_id;default:0;comment:管理员ID"`
	CategoryId   uint    `gorm:"column:category_id;default:0;comment:分类ID(单选)"`
	CategoryIds  string  `gorm:"column:category_ids;comment:分类ID(多选)"`
	Tags         string  `gorm:"column:tags;size:255;comment:标签"`
	Week         string  `gorm:"column:week;type:enum('monday','tuesday','wednesday');comment:星期(单选):monday=星期一,tuesday=星期二,wednesday=星期三"`
	Flag         string  `gorm:"column:flag;type:set('hot','index','recommend');comment:标志(多选):hot=热门,index=首页,recommend=推荐"`
	Genderdata   string  `gorm:"column:genderdata;type:enum('male','female');default:male;comment:性别(单选):male=男,female=女"`
	Hobbydata    string  `gorm:"column:hobbydata;type:set('music','reading','swimming');comment:爱好(多选):music=音乐,reading=读书,swimming=游泳"`
	Title        string  `gorm:"column:title;size:100;comment:标题"`
	Content      string  `gorm:"column:content;comment:内容"`
	Image        string  `gorm:"column:image;size:100;comment:图片"`
	Images       string  `gorm:"column:images;size:1500;comment:图片组"`
	Attachfile   string  `gorm:"column:attachfile;size:100;comment:附件"`
	Keywords     string  `gorm:"column:keywords;size:255;comment:关键字"`
	Description  string  `gorm:"column:description;size:255;comment:描述"`
	City         string  `gorm:"column:city;size:100;comment:省市"`
	Array        string  `gorm:"column:array;comment:数组:value=值"`
	Json         string  `gorm:"column:json;comment:配置:key=名称,value=值"`
	Multiplejson string  `gorm:"column:multiplejson;size:1500;comment:二维数组:title=标题,intro=介绍,author=作者,age=年龄"`
	Price        float64 `gorm:"column:price;type:decimal(10,2);default:0.00;comment:价格"`
	Views        uint    `gorm:"column:views;default:0;comment:点击"`
	Workrange    string  `gorm:"column:workrange;size:100;comment:时间区间"`
	Startdate    string  `gorm:"column:startdate;type:date;comment:开始日期"`
	Activitytime string  `gorm:"column:activitytime;type:datetime;comment:活动时间(datetime)"`
	Year         int     `gorm:"column:year;type:year;comment:年"`
	Times        string  `gorm:"column:times;type:time;comment:时间"`
	Refreshtime  int64   `gorm:"column:refreshtime;comment:刷新时间"`
	Createtime   int64   `gorm:"column:createtime;autoCreateTime;comment:创建时间"`
	Updatetime   int64   `gorm:"column:updatetime;autoUpdateTime;comment:更新时间"`
	Deletetime   int64   `gorm:"column:deletetime;comment:删除时间"`
	Weigh        int     `gorm:"column:weigh;default:0;comment:权重"`
	Switch       bool    `gorm:"column:switch;default:false;comment:开关"`
	Status       string  `gorm:"column:status;type:enum('normal','hidden');default:normal;comment:状态"`
	State        string  `gorm:"column:state;type:enum('0','1','2');default:1;comment:状态值:0=禁用,1=正常,2=推荐"`
}
