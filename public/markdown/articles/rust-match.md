首先看一段代码：
```rust
fn main() {
    let age = Some(30);
    println!("在匹配前，age是{:?}", age);
    match age {
        Some(x) =>  println!("匹配出来的age是{}", x),
        _ => ()
    }
    println!("在匹配后，age是{:?}", age);

// ======================
    let v = Some(3u8);
    match v {
        Some(3) => println!("three"),
        _ => (),
    }
 }
```
这段代码的上半部分，我认为是既用到了解构，又用到了模式匹配。
是把这两合并到了一起。
那么这种写法还可以等于：
```rust
let age = Some(18);
if let Some(age) = age {
   println!("这里永远成立");
}
```
也就是说如果match匹配里面只想关注其中一个分支的话，那么就可以简写为上面这种形式。
这里的这种语法就是·`if let` 语法。
并且有人觉得这种等效的写法是没有意义的，其实是有意义的。
比如下面这种代码：
```rust
fn main() {
    let age: Option<i32> = None;
    if let Some(age) = age {
        println!("这里永远成立");
    } else {
        println!("不成立");
    }
}
```
如果给age赋值给None，那么就不会匹配成功。就会进入else分支。

所以if let语法也可以看成是一种简写。

继续看一个代码：
```rust
fn main() {
    let x = Some(5);
    let y = 10;

    match x {
        Some(50) => println!("Got 50"),
        Some(y) => println!("Matched, y = {:?}", y),
        _ => println!("Default case, x = {:?}", x),
    }

    println!("at the end: x = {:?}, y = {:?}", x, y);
}
```
让我们看看当 match 语句运行的时候发生了什么。第一个匹配分支的模式并不匹配 x 中定义的值，所以代码继续执行。

第二个匹配分支中的模式引入了一个新变量 y，它会匹配任何 Some 中的值。因为这里的 y 在 match 表达式的作用域中，而不是之前 main 作用域中，所以这是一个新变量，不是开头声明为值 10 的那个 y。这个新的 y 绑定会匹配任何 Some 中的值，在这里是 x 中的值。因此这个 y 绑定了 x 中 Some 内部的值。这个值是 5，所以这个分支的表达式将会执行并打印出 Matched，y = 5。

如果 x 的值是 None 而不是 Some(5)，头两个分支的模式不会匹配，所以会匹配模式 _。这个分支的模式中没有引入变量 x，所以此时表达式中的 x 会是外部没有被遮蔽的 x，也就是 None。

一旦 match 表达式执行完毕，其作用域也就结束了，同理内部 y 的作用域也结束了。最后的 println! 会打印 at the end: x = Some(5), y = 10。

> 关于解构时的_

如下代码：
```rust
let mut setting_value = Some(5);
let new_setting_value = Some(10);

match (setting_value, new_setting_value) {
    (Some(_), Some(_)) => {
        println!("Can't overwrite an existing customized value");
    }
    _ => {
        setting_value = new_setting_value;
    }
}

println!("setting is {:?}", setting_value);
```
这段代码会打印出 Can't overwrite an existing customized value 接着是 setting is Some(5)。

第一个匹配分支，我们不关心里面的值，只关心元组中两个元素的类型，因此对于 Some 中的值，直接进行忽略。 剩下的形如 (Some(_),None)，(None, Some(_)), (None,None) 形式，都由第二个分支 _ 进行分配。


> 解构的时候的所有权问题

```rust
let s = Some(String::from("Hello!"));

if let Some(_s) = s {
    println!("found a string");
}

println!("{:?}", s);
```
s 是一个拥有所有权的动态字符串，在上面代码中，我们会得到一个错误，因为 s 的值会被转移给 _s，在 println! 中再次使用 s 会报错：
```rust
error[E0382]: borrow of partially moved value: `s`
 --> src/main.rs:8:22
  |
4 |     if let Some(_s) = s {
  |                 -- value partially moved here
...
8 |     println!("{:?}", s);
  |                      ^ value borrowed here after partial move
  ```
只使用下划线本身，则并不会绑定值，因为 s 没有被移动进 _：

```rust
let s = Some(String::from("Hello!"));

if let Some(_) = s {
    println!("found a string");
}

println!("{:?}", s);
found a string
Some("Hello!")
```




给出一些参考链接：
> matches! 宏的使用
> 
https://mp.weixin.qq.com/s/u1x3-c4M53W-1b6kT_rgOQ

