package logger

import (
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewDailyFileZapLogger(dir string) *zap.Logger {

	keepDays := 30 * 24 * time.Hour

	rl, err := rotatelogs.New(
		dir+"%Y-%m-%d.log",
		rotatelogs.WithLinkName(dir+".log"),       // current symlink
		rotatelogs.WithRotationTime(24*time.Hour), // har 24 soatda (kunlik)
		rotatelogs.WithMaxAge(keepDays),           // eski loglarni o‘chiradi
	)
	if err != nil {
		panic(err)
	}

	encCfg := zap.NewProductionEncoderConfig()
	encCfg.TimeKey = "ts"
	encCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	encoder := zapcore.NewJSONEncoder(encCfg)

	var ws = zapcore.AddSync(rl)

	// ws = zapcore.NewMultiWriteSyncer(ws, zapcore.AddSync(os.Stdout))

	core := zapcore.NewCore(encoder, ws, zapcore.InfoLevel)
	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}
