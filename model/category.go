/*
 * @Author: link a624669980@163.com
 * @Date: 2022-05-16 17:37:08
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-06-22 17:45:53
 * @FilePath: /CasaOS/model/category.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE

 */
package model

type ServerCategoryList struct {
	Version string         `json:"version"`
	Item    []CategoryList `json:"item"`
}
type CategoryList struct {
	Id uint `gorm:"column:id;primary_key" json:"id"`
	//CreatedAt time.Time `json:"created_at"`
	//
	//UpdatedAt time.Time `json:"updated_at"`
	Font  string `json:"font"`
	Name  string `json:"name"`
	Count uint   `json:"count"`
}
