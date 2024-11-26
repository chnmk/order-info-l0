package models

import "context"

type Consumer interface {
	Read(context.Context)
}
