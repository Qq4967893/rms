## panic
- 在程序启动的时候，如果有强依赖的服务出现故障时(例如 MySQL、Redis)，可以 panic 退出
- 在程序启动的时候，如果发现有配置明显不符合要求， 可以 panic 退出（防御编程）
- 其他情况下只要不是不可恢复的程序错误，都不应该直接 panic，而应该返回 error
- 在程序入口处，例如 gin 中间件需要使用 recover 函数，预防 panic 程序退出

## error
#### **函数使用**
使用 [pkg/errors](https://github.com/pkg/errors) 包的方法来创建和包装业务错误。
- Errorf(format, args) 和 New(message) 方法用来根据错误消息生成一个错误
- Wrap(err) 方法用来包装底层错误，增加上下文错误信息并附加调用栈； 一般用于包装第三方代码（标准库或第三方库）返回的错误
- WithMessage(err, message) 方法仅增加上下文文本信息，不附加调用栈。 如果确定错误已被 Wrap 过或不关心调用栈，可以使用此方法
- Cause(err) 方法用来获取原始错误，可以简单理解为第一次调用 Wrap(err) 方法时作为参数传入的错误
- Is(err, target) 方法用来判断错误栈里是否包含目标错误，用于匹配错误，避免直接用 `==` 去匹配，因为如果原始错误调用过 Wrap(err) 方法之后使用`==`无法正确匹配
- As(err, target) 方法用来判断错误栈里是否包含目标错误，如果包含则将错误赋值给 target
- 使用 fmt.Printf("%v", err) 或 fmt.Printf("%+v", err) 打印错误消息，前者不包含调用栈信息

#### **注意事项**
- 不要反复使用 Wrap(err) 方法 ，不然在打印错误堆栈时会显示重复的调用栈
- 调用 Wrap(err) 包装第三方库的错误，会使上层的代码与第三方包绑定死，例如有代码 errors.Wrap(sql.ErrNoRows)，则上层代码可以通过 errors.Is(err, sql.ErrNoRows) 判断具体类型，而后期如果要更换 `database/sql` 包则没有那么方便

#### **Ref**
- [Go 1.13中的错误处理](https://tonybai.com/2019/10/18/errors-handling-in-go-1-13/)
- [错误处理一：最佳实践](https://yuyy.info/?p=2003)