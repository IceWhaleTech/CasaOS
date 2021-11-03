package sort

import (
	"sort"

	"github.com/IceWhaleTech/CasaOS/model"
)

// 数据集类型, 与上一篇排序文章(多字段单独排序)比较, less字段的数据类型不再是 func(p1, p2 *Change) bool
// 而是 []func(p1, p2 *Change) bool 因为在第一个比较的值相等的情况下, 还要比较第二个值, 所以这里需要多个比较函数
type evnSorter struct {
	evn  []model.Envs
	less []lessFuncEnv
}

// sort接口方法之一(Less)
type lessFuncEnv func(p1, p2 *model.Envs) bool

// Sort 函数有两个作用
// 第一, 将参数(实际的数据集)赋值给ms对象
// 第二, 调用内置sort函数进行排序操作
func (ms *evnSorter) Sort(env []model.Envs) {
	ms.evn = env
	sort.Sort(ms)
}

// OrderedBy 函数的作用是返回一个multiSorter实例, 并将所有的实际排序函数赋值给实例的less字段,
// 上面已经为multiSorter结构体定义了Sort方法, 所以该函数的返回值可以直接调用Sort方法进行排序
// 该函数中, 为multiSorter结构体中的less字段赋值, Sort方法中又将实际数据集传入, 赋值给multiSorter的ports字段
// 一个函数, 一个方法调用过后, multiSorter实例中两个字段就已经全部被正确赋值, 可以调用系统sort函数进行排序
// 该函数也可看作是一个工厂方法, 用来生成less字段已经被赋值的multiSorter实例
func EnvSort(less ...lessFuncEnv) *evnSorter {
	return &evnSorter{
		less: less,
	}
}

// Len 为sort接口方法之一
func (ms *evnSorter) Len() int {
	return len(ms.evn)
}

// Swap 为sort接口方法之一
func (ms *evnSorter) Swap(i, j int) {
	ms.evn[i], ms.evn[j] = ms.evn[j], ms.evn[i]
}

// Less 为sort接口方法之一
func (ms *evnSorter) Less(i, j int) bool {
	temp := ms.evn
	p, q := &temp[i], &temp[j]
	// Try all but the last comparison.
	var k int
	// 由于可能有多个需要排序的字段, 也就对应了多个less函数, 当第一个字段的值相等时,
	// 需要依次尝试比对后续其他字段的值得大小, 所以这里需要获取比较函数的长度, 以便遍历比较
	for k = 0; k < len(ms.less)-1; k++ {
		// 提取比较函数, 将函数赋值到新的变量中以便调用
		less := ms.less[k]
		switch {
		case less(p, q):
			// 如果 p < q, 返回值为true, 不存在两个值相等需要比较后续字段的情况, 所以这里直接返回
			// 如果 p > q, 返回值为false, 则调到下一个case中处理
			return true
		case less(q, p):
			// 如果 p > q, 返回值为false, 不存在两个值相等需要比较后续字段的情况, 所以这里直接返回
			return false
		}
		// 如果代码走到这里, 说明ms.less[k]函数比较后 p == q; 重新开始下一次循环, 更换到下一个比较函数处理
		continue
	}
	// 如果代码走到这里, 说明所有的比较函数执行过后, 所有比较的值都相等
	// 直接返回最后一次的比较结果数据即可
	return ms.less[k](p, q)
}
