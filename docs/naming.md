# Naming
## Package naming
当命名包时，请按下面规则选择一个名称：

- 全部小写。没有大写或下划线。
- 大多数使用命名导入的情况下，不需要重命名。
- 简短而简洁。请记住，在每个使用的地方都完整标识了该名称。
- 不用复数。例如net/url，而不是net/urls。
- 不要用“common”，“util”，“shared”或“lib”。这些是不好的，信息量不足的名称。

另请参阅 [Go 包命名规则](https://go.dev/blog/package-names) 和 [Go 包样式指南](https://rakyll.org/style-packages/).

## File naming

## Function naming
我们遵循 Go 社区关于使用 [MixedCaps 作为函数名](https://go.dev/doc/effective_go#mixed-caps) 的约定。有一个例外，为了对相关的测试用例进行分组，函数名可能包含下划线，如：TestMyFunction_WhatIsBeingTested.

### 业务逻辑方法命名
Controller:
- IsEmailExisted()
- IsUsernameExisted()
- GetInfo()
- GetAll()
- GetList()
- Add()
- Update()
- Delete()
- buildXxx()

Service:
命名原则与 Controller 层类似，只不过 Service 主要承担了业务逻辑的责任，会包含更多的不同条件的查询方法，以及调用其他的 Service 或 Repository 方法进行整合数据。
- IsIdExisted()
- IsNameExisted()
- GetInfoById()
- GetInfoByName()
- GetListByRoleId()
- GetXxxList()
- GetXxxMap()


Repository: 基础的数据库操作方法，不应该包含过多的业务逻辑处理。
> 方法的顺序（或者叫分组）最好保持一定的规范，方便查找和添加方法有固定的逻辑


## Error naming
对于存储为全局变量的错误值， 根据是否导出，使用前缀 Err 或 err。 请看指南 对于未导出的顶层常量和变量，使用_作为前缀。
``` go
var (
  // 导出以下两个错误，以便此包的用户可以将它们与 errors.Is 进行匹配。
  ErrBrokenLink = errors.New("link is broken")
  ErrCouldNotOpen = errors.New("could not open")

  // 这个错误没有被导出，因为我们不想让它成为我们公共 API 的一部分。 我们可能仍然在带有错误的包内使用它。
  errNotFound = errors.New("not found")
)
```
对于自定义错误类型，请改用后缀 `Error`。
``` go
// 同样，这个错误被导出，以便这个包的用户可以将它与 errors.As 匹配。
type NotFoundError struct {
  File string
}

func (e *NotFoundError) Error() string {
  return fmt.Sprintf("file %q not found", e.File)
}

// 并且这个错误没有被导出，因为我们不想让它成为公共 API 的一部分。 我们仍然可以在带有 errors.As 的包中使用它。
type resolveError struct {
  Path string
}

func (e *resolveError) Error() string {
  return fmt.Sprintf("resolve %q", e.Path)
}
```
