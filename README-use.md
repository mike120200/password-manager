# password_manager usager
<p align="center">
  <a href="#zh">中文</a> | <a href="#en">English</a>
</p>


## <a id="zh"></a>📌 中文
<details open>
<summary>展开/折叠</summary>



### 添加密码

#### 简介：存储账号、密码、此账号密码使用的平台（可选填，直接按enter即可）

#### 使用方法：

```sh
pm add
Enter Key or account : key_name
Enter password : ******
Enter platform(optional) : 
```

```sh
pm add <key_name>
Enter password : ******
Enter platform(optional) : 
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

### 根据平台信息查找

#### 简介：输入平台信息的关键词进行查找

```sh
pm pla
Enter platform : 
```

```sh
pm pla <platform>
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

### 从备份文件恢复数据

#### 简介：从备份文件恢复数据

#### 使用方法：

```sh
pm restore
```

</details>

## <a id="en"></a>📌 English
<details open>
<summary>展开/折叠</summary>
### Add a Password

#### Description: Store the account, password, and the platform used for this account password (optional, press enter)

#### Usage:

```sh
pm add
Enter Key or account : key_name
Enter password : ******
Enter platform(optional) :
```

```sh
pm add <key_name>
Enter password : ******
Enter platform(optional) : 
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

### Search by Platform Information

#### Introduction: Enter keywords related to the platform information to perform a search.

```sh
pm pla
Enter platform : 
```

```sh
pm pla <platform>
```

---

### 

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

### Restore form backup file

#### Description:Restore datas from backup files

#### Usage:

```sh
pm restore
```



</details>