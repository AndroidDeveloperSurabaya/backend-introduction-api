package data

import (
	"errors"
	"fmt"
	"log"
	"time"

	"amalhanaja.com/user-service/domain"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" //
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// NewUser ..
type NewUser struct {
	Password string `json:"password"`
	FullName string `json:"fullName"`
}

type userDataRepository struct {
	cache *redis.Client
	db    *gorm.DB
}

// User ...
type User struct {
	UUID            []byte    `json:"uuid" gorm:"type:uuid;primary_key;not null;default:uuid_generate_v4()"`
	FullName        string    `json:"fullName" gorm:"type:varchar(200);not null"`
	Password        string    `json:"password" gorm:"type:varchar(128);not null"`
	IsActive        bool      `json:"isActive" gorm:"not null;default:false"`
	Email           string    `json:"email" gorm:"unique;not null"`
	IsEmailVerified bool      `json:"isEmailVerified" gorm:"not null;default:false"`
	CreatedAt       time.Time `json:"createdAt" gorm:"not null"`
	LastUpdatedAt   time.Time `json:"lastUpdatedAt" gorm:"not null"`
}

var instance *userDataRepository

// GetInstance ...
func GetInstance() domain.UserRepository {
	if instance == nil {
		db, err := gorm.Open("postgres", "host=localhost port=5432 user=amalhanaja dbname=amalhanaja_user password=tanggallahir1998 sslmode=disable")
		if err != nil {
			log.Panic(err)
			return nil
		}
		instance = &userDataRepository{
			cache: redis.NewClient(&redis.Options{
				Addr:     "localhost:6379",
				Password: "",
				DB:       0,
			}),
			db: db,
		}
		instance.migrate(&User{})
	}
	return instance
}

func (repo *userDataRepository) migrate(tables ...interface{}) {
	if repo.db.HasTable(&User{}) {
		repo.db.AutoMigrate(&User{})
	} else {
		repo.db.CreateTable(&User{})
	}
}

func (repo *userDataRepository) StoreUser(newUser domain.NewUser) (string, error) {
	now := time.Now()
	if u, err := uuid.NewV4(); err == nil {
		claims := &jwt.StandardClaims{
			IssuedAt:  now.Unix(),
			ExpiresAt: (now.Unix() + ExpirationTimeInMillis),
			Id:        u.String(),
			Subject:   newUser.Email,
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		if signedToken, err := token.SignedString([]byte(JwtSecretKey)); err == nil {
			passwordHash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
			if err != nil {
				return "", err
			}
			user := &User{
				FullName:        newUser.FullName,
				Email:           newUser.Email,
				Password:        string(passwordHash),
				IsEmailVerified: false,
				IsActive:        false,
				CreatedAt:       now.UTC(),
				LastUpdatedAt:   now.UTC(),
			}
			log.Println(repo.db.NewRecord(user))
			if repo.db.NewRecord(user) {
				if err := repo.db.Create(user).Error; err != nil {
					return "", err
				}
				activationURL := fmt.Sprintf("http://192.168.1.15:8080/activate?activationToken=%s", signedToken)
				mail := NewMail(newUser.Email, "Verification", activationURL)
				if err := mail.SendMessage(); err != nil {
					log.Panic(errors.New("Can't Send an Email"), err)
					return "", err
				}
				return signedToken, nil
			}
			log.Panic(errors.New("Can't Make Json Binary"), err)
			return "", err
		}
		log.Panic(errors.New("Error Creating SignedToken"), err)
	} else {
		log.Panic(err)
		return "", err
	}
	return "", errors.New("Unknown Error")
}

func (repo *userDataRepository) GetUserByEmail(email string) (*domain.User, error) {
	user := &User{}
	if err := repo.db.First(user, &User{Email: email}).Error; err != nil {
		return nil, err
	}
	log.Println(user)
	return &domain.User{
		Email:           user.Email,
		FullName:        user.FullName,
		IsEmailVerified: user.IsEmailVerified,
		IsActive:        user.IsActive,
		CreatedAt:       user.CreatedAt,
		LastModifiedAt:  user.LastUpdatedAt,
	}, nil
}

func (repo *userDataRepository) GetUserByUUID(uuid string) (*domain.User, error) {
	user := &User{}
	if err := repo.db.First(user, &User{UUID: []byte(uuid)}).Error; err != nil {
		return nil, err
	}
	log.Println(user)
	return &domain.User{
		Email:           user.Email,
		FullName:        user.FullName,
		IsEmailVerified: user.IsEmailVerified,
		IsActive:        user.IsActive,
		CreatedAt:       user.CreatedAt,
		LastModifiedAt:  user.LastUpdatedAt,
	}, nil
}

func (repo *userDataRepository) GetUserByToken(token string) (*domain.User, error) {
	parsedToken, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected Signing Method")
		}
		return []byte(JwtSecretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		subject := claims["sub"].(string)
		now := time.Now()
		if newUser, err := repo.GetUserByEmail(subject); err == nil {
			return &domain.User{
				Email:          newUser.Email,
				FullName:       newUser.FullName,
				CreatedAt:      now.UTC(),
				LastModifiedAt: now.UTC(),
			}, nil
		} else {
			if user, err := repo.GetUserByUUID(subject); err == nil {
				return &domain.User{
					Email:          user.Email,
					FullName:       user.FullName,
					CreatedAt:      now.UTC(),
					LastModifiedAt: now.UTC(),
				}, nil
			}
		}

		return nil, errors.New("Error claims")
	}
	return nil, errors.New("Invalid Token")
}
func (repo *userDataRepository) UpdateUser(userDomain *domain.User) (*domain.User, error) {
	user := &User{}
	if err := repo.db.First(user, &User{Email: userDomain.Email}).Error; err != nil {
		return nil, err
	}
	user.Email = userDomain.Email
	user.FullName = userDomain.FullName
	user.IsActive = userDomain.IsActive
	user.IsEmailVerified = userDomain.IsEmailVerified
	user.LastUpdatedAt = time.Now().UTC()
	if err := repo.db.Save(user).Error; err != nil {
		return nil, err
	}
	return userDomain, nil
}

func (repo *userDataRepository) CreateJWTToken(email string, password string) (*domain.JwtToken, error) {
	user := &User{}
	now := time.Now()
	if err := repo.db.First(user, &User{Email: email}).Error; err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err
	}

	if !user.IsActive {
		return nil, errors.New("User Not Active")
	}
	uid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	claims := &jwt.StandardClaims{
		IssuedAt:  now.Unix(),
		ExpiresAt: (now.Unix() + ExpirationTimeInMillis),
		Id:        uid.String(),
		Subject:   string(user.UUID[:]),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if signedToken, err := token.SignedString([]byte(JwtSecretKey)); err == nil {
		fmt.Println(signedToken)
		return &domain.JwtToken{
			Token:     signedToken,
			ExpiresIn: ExpirationTimeInMillis,
		}, nil
	}
	return nil, err

}
