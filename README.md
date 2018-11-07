<h1 align="center">给不了你梦中情人，至少还有硬盘女神：hardseedGO</h1>

https://github.com/yangyangwithgnu/hardseed/
这里太久没有更新了。虽然还能正常运行，但是下载种子文件的时候草榴出现错误

但是原著的程序是用C++写的，下载才疏学浅，也懒得改C系的程序，所以用GOlang重新实现以下

目前还在开发中，暂时没有提供各种版本的支持。
在开发接近尾声的时候，会发布所有系统的支持包。
因为GOlang的特性，很容易实现跨平台的支持。所以各位放心。
<br />

todo
<br />
1. ~~草榴支持~~
2. ~~1024支持~~
3. ~~不重复下载种子和图片~~ 如果想重置，请删除data.db
4. 自动提交种子到下载器
<br />
目前支持
<br />

- xp_asia_mosaiched
- xp_asia_non_mosaiched
- aicheng_asia_mosaiched
- aicheng_asia_non_mosaicked
- chaoliu_asia_mosaiched
- chaoliu_asia_non_mosaiched

<br />
6个class
<br />
可以一次定义多个class，默认会按照顺序自行下载。
<br />
其他设置请参考config.yaml自行修改。
<br />