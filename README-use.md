# password_manager usager
<p align="center">
  <a href="#zh">中文</a> | <a href="#en">English</a>
</p>

## <a id="zh"></a>📌 中文
<details open>
<summary>展开/折叠</summary>
### 添加密码

#### 简介：存储一个密码对

#### 使用方法：

```sh
pm add
Enter Key : key_name
Enter password : ******
```

```sh
pm add <key_name>
Enter password : ******
```

---

### 查找密码

#### 简介：根据唯一标识查找密码

#### 使用方法：

```sh
pm query
Enter Key : key_name
```

```sh
pm query <key_name>
```

---

### 列出所有已存密码

#### 简介：列出所有已存密码

#### 使用方法：

```
pm list
```

---

### 更新密码

#### 简介：更新密码

#### 使用方法：

```sh
pm update
Enter Key : key_name
Enter password : *******
```

```sh
pm update <key_name>
Enter password : *******
```

---

### 删除密码

#### 简介：删除密码

#### 使用方法：

```sh
pm del
Enter Key : key_name
```

```sh
pm del <key_name>
```

---

### 备份密码

#### 简介：每隔一段时间会自动备份一次，也可手动备份

#### 使用方法：

```shell
pm backup
```



</details>

## <a id="en"></a>📌 English
<details open>
<summary>展开/折叠</summary>
### Add a Password

#### Description: Store a password key-value pair.

#### Usage:

```sh
pm add
Enter Key: key_name
Enter password: ******
```

```sh
pm add <key_name>
Enter password: ******
```

---

### **Query a Password**

#### **Description: Retrieve a stored password using its unique identifier.**

#### Usage:

```sh
pm query
Enter Key: key_name
```

```sh
pm query <key_name>
```

---

### **List All Stored Passwords**

#### **Description: Display all stored passwords.**

#### Usage:

```sh
pm list
```

---

### **Update a Password**

#### **Description: Modify an existing stored password.**

#### Usage:

```sh
pm update
Enter Key: key_name
Enter password: *******
```

```sh
pm update <key_name>
Enter password: *******
```

---

### **Delete a Password**

#### **Description: Remove a stored password.**

#### Usage:

```sh
pm del
Enter Key: key_name
```

```sh
pm del <key_name>
```

---

### **Backup Passwords**

#### **Description: Passwords are automatically backed up at regular intervals, but manual backups are also possible.**

#### Usage:

```sh
pm backup
```



</details>