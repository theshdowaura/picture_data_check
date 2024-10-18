package test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

// 定义结构体
type Time struct {
	Start *int `json:"start"`
	Last  *int `json:"last"`
}

type Course struct {
	Course      string `json:"course"`
	Place       string `json:"place"`
	Campus      string `json:"campus"`
	Teacher     string `json:"teacher"`
	WeeksText   string `json:"weeks_text"`
	WeekDay     string `json:"week_day"`
	WeekDayText string `json:"week_day_text"`
	TimeText    string `json:"time_text"`
	WeeksArr    []int  `json:"weeks_arr"`
	Time        Time   `json:"time"`
	Color       string `json:"color"`
	Section     string `json:"section"`
}

type Schedule struct {
	ClassName  string   `json:"class_name"`
	Username   string   `json:"username"`
	CourseList []Course `json:"course_list"`
}

func TestParseSchedule(t *testing.T) {
	// 模拟的 JSON 数据
	data := `
{'class_name': '信安224', 'username': '姜浩杰', 'course_list': [{'course': '信息安全项目实战', 'place': '教学楼E105', 'campus': '大学城校区', 'teacher': '李建新', 'weeks_text': '9-11周', 'week_day': '1', 'week_day_text': '星期一', 'time_text': '星期一 1-4节', 'weeks_arr': [9, 10, 11], 'time': {'start': None, 'last': None}, 'color': 'green', 'section': '1-4'}, {'course': '现代密码学', 'place': '工业互联网大楼0508', 'campus': '大学城校区', 'teacher': '张静', 'weeks_text': '1-2周', 'week_day': '1', 'week_day_text': '星期一', 'time_text': '星期一 5-8节', 'weeks_arr': [1, 2], 'time': {'start': None, 'last': None}, 'color': 'yellow', 'section': '5-8'}, {'course': '综合项目实践', 'place': '教学楼A111', 'campus': '大学城校区', 'teacher': '吴昊', 'weeks_text': '4-5周', 'week_day': '1', 'week_day_text': '星期一', 'time_text': '星期一 5-8节', 'weeks_arr': [4, 5], 'time': {'start': None, 'last': None}, 'color': 'yellow', 'section': '5-8'}, {'course': '工业数据安全', 'place': '实训楼SA518', 'campus': '大学城校区', 'teacher': '乔虹', 'weeks_text': '7-8周', 'week_day': '1', 'week_day_text': '星期一', 'time_text': '星期一 5-8节', 'weeks_arr': [7, 8], 'time': {'start': None, 'last': None}, 'color': 'yellow', 'section': '5-8'}, {'course': '现代密码学', 'place': '工业互联网大楼0508', 'campus': '大学城校区', 'teacher': '张静', 'weeks_text': '1-2周', 'week_day': '2', 'week_day_text': '星期二', 'time_text': '星期二 1-4节', 'weeks_arr': [1, 2], 'time': {'start': None, 'last': None}, 'color': 'yellow', 'section': '1-4'}, {'course': '综合项目实践', 'place': '教学楼A111', 'campus': '大学城校区', 'teacher': '吴昊', 'weeks_text': '4周', 'week_day': '2', 'week_day_text': '星期二', 'time_text': '星期二 1-4节', 'weeks_arr': [4], 'time': {'start': None, 'last': None}, 'color': 'yellow', 'section': '1-4'}, {'course': '综合项目实践', 'place': '教学楼A304', 'campus': '大学城校区', 'teacher': '吴昊', 'weeks_text': '6周', 'week_day': '2', 'week_day_text': '星期二', 'time_text': '星期二 1-4节', 'weeks_arr': [6], 'time': {'start': None, 'last': None}, 'color': 'yellow', 'section': '1-4'}, {'course': '工业数据安全', 'place': '实训楼SA518', 'campus': '大学城校区', 'teacher': '乔虹', 'weeks_text': '7-8周', 'week_day': '2', 'week_day_text': '星期二', 'time_text': '星期二 1-4节', 'weeks_arr': [7, 8], 'time': {'start': None, 'last': None}, 'color': 'yellow', 'section': '1-4'}, {'course': '信息安全项目实战', 'place': '教学楼E105', 'campus': '大学城校区', 'teacher': '李建新', 'weeks_text': '9-11周', 'week_day': '2', 'week_day_text': '星期二', 'time_text': '星期二 1-4节', 'weeks_arr': [9, 10, 11], 'time': {'start': None, 'last': None}, 'color': 'green', 'section': '1-4'}, {'course': '综合项目实践', 'place': '教学楼A304', 'campus': '大学城校区', 'teacher': '吴昊', 'weeks_text': '4-6周(双)', 'week_day': '3', 'week_day_text': '星期三', 'time_text': '星期三 1-4节', 'weeks_arr': [4, 6], 'time': {'start': None, 'last': None}, 'color': 'yellow', 'section': '1-4'}, {'course': '信息安全项目实战', 'place': '教学楼E105', 'campus': '大学城校区', 'teacher': '李建新', 'weeks_text': '9-11周', 'week_day': '3', 'week_day_text': '星期三', 'time_text': '星期三 1-4节', 'weeks_arr': [9, 10, 11], 'time': {'start': None, 'last': None}, 'color': 'green', 'section': '1-4'}, {'course': '现代密码学', 'place': '工业互联网大楼0508', 'campus': '大学城校区', 'teacher': '张静', 'weeks_text': '1-2周', 'week_day': '3', 'week_day_text': '星期三', 'time_text': '星期三 1-8节', 'weeks_arr': [1, 2], 'time': {'start': None, 'last': None}, 'color': 'yellow', 'section': '1-8'}, {'course': '工业数据安全', 'place': '实训楼SA518', 'campus': '大学城校区', 'teacher': '乔虹', 'weeks_text': '7-8周', 'week_day': '3', 'week_day_text': '星期三', 'time_text': '星期三 5-8节', 'weeks_arr': [7, 8], 'time': {'start': None, 'last': None}, 'color': 'yellow', 'section': '5-8'}, {'course': '现代密码学', 'place': '工业互联网大楼0508', 'campus': '大学城校区', 'teacher': '张静', 'weeks_text': '1-2周', 'week_day': '4', 'week_day_text': '星期四', 'time_text': '星期四 1-4节', 'weeks_arr': [1, 2], 'time': {'start': None, 'last': None}, 'color': 'yellow', 'section': '1-4'}, {'course': '信息安全项目实战', 'place': '教学楼E105', 'campus': '大学城校区', 'teacher': '李建新', 'weeks_text': '9-11周', 'week_day': '4', 'week_day_text': '星期四', 'time_text': '星期四 1-4节', 'weeks_arr': [9, 10, 11], 'time': {'start': None, 'last': None}, 'color': 'green', 'section': '1-4'}, {'course': '工业数据安全', 'place': '实训楼SA518', 'campus': '大学城校区', 'teacher': '乔虹', 'weeks_text': '7-8周', 'week_day': '4', 'week_day_text': '星期四', 'time_text': '星期四 1-8节', 'weeks_arr': [7, 8], 'time': {'start': None, 'last': None}, 'color': 'yellow', 'section': '1-8'}, {'course': '综合项目实践', 'place': '教学楼A111', 'campus': '大学城校区', 'teacher': '吴昊', 'weeks_text': '3-4周,6周', 'week_day': '4', 'week_day_text': '星期四', 'time_text': '星期四 5-8节', 'weeks_arr': [3, 4, 6], 'time': {'start': None, 'last': None}, 'color': 'yellow', 'section': '5-8'}, {'course': '现代密码学', 'place': '工业互联网大楼0508', 'campus': '大学城校区', 'teacher': '张静', 'weeks_text': '1-2周', 'week_day': '5', 'week_day_text': '星期五', 'time_text': '星期五 1-4节', 'weeks_arr': [1, 2], 'time': {'start': None, 'last': None}, 'color': 'yellow', 'section': '1-4'}, {'course': '综合项目实践', 'place': '教学楼A209', 'campus': '大学城校区', 'teacher': '吴昊', 'weeks_text': '4周', 'week_day': '5', 'week_day_text': '星期五', 'time_text': '星期五 1-4节', 'weeks_arr': [4], 'time': {'start': None, 'last': None}, 'color': 'yellow', 'section': '1-4'}, {'course': '工业数据安全', 'place': '实训楼SA518', 'campus': '大学城校区', 'teacher': '乔虹', 'weeks_text': '7-8周', 'week_day': '5', 'week_day_text': '星期五', 'time_text': '星期五 1-4节', 'weeks_arr': [7, 8], 'time': {'start': None, 'last': None}, 'color': 'yellow', 'section': '1-4'}, {'course': '信息安全项目实战', 'place': '教学楼E105', 'campus': '大学城校区', 'teacher': '李建新', 'weeks_text': '9-11周', 'week_day': '5', 'week_day_text': '星期五', 'time_text': '星期五 1-4节', 'weeks_arr': [9, 10, 11], 'time': {'start': None, 'last': None}, 'color': 'green', 'section': '1-4'}, {'course': '综合项目实践', 'place': '教学楼A304', 'campus': '大学城校区', 'teacher': '吴昊', 'weeks_text': '3-4周', 'week_day': '5', 'week_day_text': '星期五', 'time_text': '星期五 5-8节', 'weeks_arr': [3, 4], 'time': {'start': None, 'last': None}, 'color': 'yellow', 'section': '5-8'}, {'course': '综合项目实践', 'place': '教学楼A111', 'campus': '大学城校区', 'teacher': '吴昊', 'weeks_text': '2周,6周', 'week_day': '6', 'week_day_text': '星期六', 'time_text': '星期六 5-8节', 'weeks_arr': [2, 6], 'time': {'start': None, 'last': None}, 'color': 'yellow', 'section': '5-8'}, {'course': '综合项目实践', 'place': '教学楼A304', 'campus': '大学城校区', 'teacher': '吴昊', 'weeks_text': '4周', 'week_day': '7', 'week_day_text': '星期日', 'time_text': '星期日 5-8节', 'weeks_arr': [4], 'time': {'start': None, 'last': None}, 'color': 'yellow', 'section': '5-8'}]}
`

	// 替换单引号为双引号，替换 None 为 null
	data = strings.ReplaceAll(data, "'", "\"")
	data = strings.ReplaceAll(data, "None", "null")

	var schedule Schedule

	// 解析 JSON 数据
	err := json.Unmarshal([]byte(data), &schedule)
	if err != nil {
		t.Fatalf("解析 JSON 失败: %v", err)
	}

	// 打印解析后的数据
	fmt.Printf("班级: %s\n", schedule.ClassName)
	fmt.Printf("学生: %s\n", schedule.Username)
	fmt.Println("课程列表:")
	for _, course := range schedule.CourseList {
		fmt.Printf("课程名称: %s\n", course.Course)
		fmt.Printf("地点: %s, 校区: %s\n", course.Place, course.Campus)
		fmt.Printf("老师: %s\n", course.Teacher)
		fmt.Printf("周次: %s, 星期: %s (%s)\n", course.WeeksText, course.WeekDay, course.WeekDayText)
		fmt.Printf("时间: %s\n", course.TimeText)
		fmt.Printf("颜色: %s, 节次: %s\n", course.Color, course.Section)
		fmt.Println("-------------")
	}

	// 断言解析是否正确
	if schedule.ClassName != "信安224" {
		t.Errorf("期望的班级名是 '信安224', 但是得到了 '%s'", schedule.ClassName)
	}
	if schedule.Username != "姜浩杰" {
		t.Errorf("期望的用户名是 '姜浩杰', 但是得到了 '%s'", schedule.Username)
	}
	if len(schedule.CourseList) == 0 {
		t.Fatalf("课程列表不能为空")
	}
}
