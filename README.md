# password_manager installation
<p align="center">
  <a href="#zh">中文</a> | <a href="#en">English</a>
</p>

## <a id="zh"></a>📌 中文
<details open>
<summary>展开/折叠</summary>

`Password Manager` 是一个安全的密码存储和管理工具。

### 📌 主要功能：
✅ 安全存储密码  
✅ AES 加密保护数据  
✅ 命令行管理密码
✅ 备份数据

---

### 🚀 安装方法（Mac & Windows）
你可以将编译后的 `pm` 文件获取到本地电脑，并将其添加到环境变量，以便在命令行中随时使用。

---

### **💻 MacOS 安装**
1. **下载 `pm` 可执行文件**

2. 创建pm_dir目录（其他命名也可，配置环境变量的时候对应上即可）

   ```sh
   mkdir pm_dir
   ```

3. 将pm移动至`pm_dir`

4. 配置环境变量

    + 设置 Zsh 终端的环境变量，`path`替换成真实的路径

      ```sh
      echo 'export PATH="/path/pm_dir:$PATH"' >> ~/.zshrc
      source ~/.zshrc
      ```

    + 设置bash环境变量，`path`替换成真实的路径

      ```sh
      echo 'export PATH="/path/pm_dir:$PATH"' >> ~/.bash_profile
      source ~/.bash_profile
      ```

5. **检查是否安装成功**

    ```sh
    pm --help
    ```

### **🖥️ Windows 安装**

1. **下载** pm.exe **可执行文件**

 	2. **创建** pm_dir **文件夹**（其他命名也可，配置环境变量的时候对应上即可）
 	3. 将pm移动至`pm_dir`
 	4. 配置环境变量
     + 打开 **控制面板 → 系统 → 高级系统设置 → 环境变量**
     + 在 **系统变量** 或 **用户变量** 里找到 Path
     + 点击 **编辑**，新增 文件夹`pm_dir`的绝对路径
     + **保存后重启 CMD**

5. **检查是否安装成功**

   ```sh
   pm --help
   ```



</details>

---
<a id="en"></a>📌 English
<details open>  
<summary>Expand/Collapse</summary>  

Password Manager is a secure password storage and management tool.

📌 Main Features:

✅ Securely store passwords
✅ AES encryption for data protection
✅ Manage passwords via the command line
✅ Backup stored data

🚀 Installation Guide (Mac & Windows)

You can download the compiled pm executable file to your local computer and add it to the system environment variables for convenient command-line usage.

💻 MacOS Installation
	1.	Download the pm executable file

	2.	Create a pm_dir directory (you can use any name, just ensure it matches in the environment variable configuration)
	
	```sh
	mkdir pm_dir  
	```

3.  **Move** pm **into** pm_dir

4. **Configure environment variables**

   + *For Zsh terminal**, replace path with th*e actual directory path

     ```sh
     echo 'export PATH="/path/pm_dir:$PATH"' >> ~/.zshrc  
     source ~/.zshrc 
     ```

   + **For Bash terminal**, replace path with the actual directory path

     ```sh
     echo 'export PATH="/path/pm_dir:$PATH"' >> ~/.bash_profile  
     source ~/.bash_profile 
     ```


5. **Verify installation**

   ```sh
   pm --help 
   ```

   

</details>
