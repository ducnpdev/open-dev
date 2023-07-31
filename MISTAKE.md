# mistake

## transaction of database
Mình hy vọng sau khi đọc bài này bạn có thể tránh những lỗi này trong dự án của chính mình. Chỉ là một lỗi nhỏ có thể khiến bạn phải hối hận.

Sẽ đối mặt là một transaction database đơn giản. Hầu như tất cả các develop điều đã làm với transaction một lần, nếu vẫn thắc mắc thì có thể tham khảo thêm.

Trong service sẽ làm chức năng quản lý user trong database, get user bằng id sau đó update fullname

Thứ tự các bước như sau:
* start 1 transaction
* lấy thông tin user bằng id
* Xác nhận user vẫn còn trạng thái active
* nếu hợp lệ thì update fullname
* commit transaction
* nếu có gì không hợp lệ hoặc error thì revert transaction

| Mình sẽ dùng gorm để truy cập database

- đoạn code repository get user băng id:
```go
func (r *userService) GetUserById(tx *gorm.DB, id int) (model.User, error) {
	var (
		user model.User
		err  error
	)
	err = tx.Table(user.Table()).Where("id = ?", id).First(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}
```

- đoạn code logic:
```go
func (s *service) UpdateUserId(ctx context.Context, id int) (error) {
    
	tx := r.gormDB.Begin()

    defer func() {
        if err != nil {
            tx.Rollback()
        }
    }()
    
    user, err := s.userService.GetUserById(tx, id)
    if err != nil {
        return  err
    }

    if user.Status != "active" {
        return  nil
    }
	
    user.UpdatedAt = time.Now()

    err = s.userService.UpdateName(tx, name)
    if err != nil {
        return  err
    }

    err = tx.Commit()
    if err != nil {
        return  err
    }

    return nil
}
```

### Vấn Đề
* tất cả các request từ client sẽ bị time-out, database sẽ thêm vào rất nhiều connection
* không có bất cứ 1 request nào có thể hoàn thành.

**Tại Sao**
vấn đề xảy ra đó là khi user có status khác `active` điều này có nghĩa là sẽ return nil và transaction sẽ không release lock trên record của user đó.

mặc dù trong function `defer` có gọi `rollback()` chỉ xảy ra khi có error, và user khác active sẽ return nil và transaction sẽ bị lock cho đến khi transaction bị timeout.

**Fix**

Để fix vấn đề trên là cực kì đơn giản, thêm `rollback()`.
```go
if user.Status != "active" {
    tx.Rollback()
    return  nil
}
```

còn nếu bạn vẫn muốn `rollback()` trong `defer` thì sẽ return error
```go
var ErrUserNotActive = errors.New("user not active")
func (s *service) UpdateUserId(ctx context.Context, id int) (error) {
    
	tx := r.gormDB.Begin()

    defer func() {
        if err != nil {
            tx.Rollback()
        }
    }()
    
    user, err := s.userService.GetUserById(tx, id)
    if err != nil {
        return  err
    }

    if user.Status != "active" {
        return ErrUserNotActive
    }
	
    user.UpdatedAt = time.Now()

    err = s.userService.UpdateName(tx, name)
    if err != nil {
        return  err
    }

    err = tx.Commit()
    if err != nil {
        return  err
    }

    return nil
}
```
