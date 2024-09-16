/*
usecase：アプリケーションのビジネスロジックを実装

具体的なデータの操作やビジネスルールの適用など、アプリケーションの核となる機能を担当。
コントローラから受け取ったデータを元に処理を行い、その結果をコントローラに返す。
*/
package usecase

import (
	"backend/model"
	"backend/repository"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type IUserUsecase interface {
	SignUp(user model.User) (model.UserResponse, error)
	LogIn(user model.User) (string, error)
}

type userUsecase struct {
	ur repository.IUserRepository
}

func NewUserUsecase(ur repository.IUserRepository) IUserUsecase {
	return &userUsecase{ur}
}

func (uu *userUsecase) SignUp(user model.User) (model.UserResponse, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return model.UserResponse{}, err
	}
	newUser := model.User{Email: user.Email, Password: string(hash)}
	if err := uu.ur.CreateUser(&newUser); err != nil {
		return model.UserResponse{}, err
	}
	resUser := model.UserResponse{
		ID:    newUser.ID,
		Email: newUser.Email,
	}
	return resUser, nil
}

func (uu *userUsecase) LogIn(user model.User) (string, error) {
	storedUser := model.User{}
	if err := uu.ur.GetUserByEmail(&storedUser, user.Email); err != nil {
		return "", err
	}
	err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": storedUser.ID,
		"exp":     time.Now().Add(time.Hour * 12).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

/*
Usecase: ビジネスロジックを実装し、アプリケーションの具体的な操作やビジネスルールの適用を担当します。通常、ユースケースは特定のアクションや操作に焦点を当てています。

Service: 特定の機能や処理を提供し、複数のユースケースや他のサービスから再利用されることが多いです。サービスは、より汎用的な機能を提供することが目的です。

Repository: データベース操作を担当し、データの永続化や取得を行います。データソースへのアクセスを抽象化し、ビジネスロジックから分離します。

Model: アプリケーションのデータ構造を定義します。データベースのエンティティやAPIのリクエスト/レスポンスの形式を表現します。

Controller: クライアントからのリクエストを受け取り、適切なユースケースを呼び出してレスポンスを返します。プレゼンテーション層の一部として機能します。

Router: ルーティングを設定し、特定のURLパスに対して適切なコントローラを呼び出します。リクエストの振り分けを行います。

Validator: 入力データの検証を行い、データの整合性や形式をチェックします。バリデーションエラーがある場合は、適切なエラーメッセージを返します。

Middleware: リクエストがルーターに到達する前に実行され、認証、ログ記録、エラーハンドリングなどの共通の機能を提供します。
*/
