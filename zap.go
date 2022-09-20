package gutils

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 基于go.uber.org/zap 源
// 参考原作者https://github.com/SliverHorn
// https://github.com/znToast 改

/*
	初始化日志

@Parameter dir     日志输出目录
@Parameter console 是否在控制台打印
@return *zap.Logger	日志对象
*/
func InitZap(dir string, con bool) (logger *zap.Logger) {
	log.Println("初始化日志...................")
	// 判断log文件夹是否存在
	if ok, _ := PathExist(dir); !ok {
		log.Printf("创建新文件夹:%v\n", dir)
		os.Mkdir(dir, os.ModePerm)
	}
	//new 一个zap
	z := New_zap(dir, con)
	core := z.GetZapCores()
	logger = zap.New(zapcore.NewTee(core...))
	// 行打印
	logger = logger.WithOptions(zap.AddCaller())
	zap.ReplaceGlobals(logger)
	return logger
}

// 适合自己的默认配置
func New_zap(dir string, console bool) *_zap {
	return &_zap{
		Level:         zapcore.InfoLevel,                  //info级别
		MaxAge:        30,                                 // 30天
		LevelEncoder:  zapcore.LowercaseColorLevelEncoder, //小写编码器带颜色
		StacktraceKey: "stacktrace",                       //栈名
		directory:     dir,                                //输出目录
		Console:       console,                            //是否在控制台打印
	}
}

type _zap struct {
	Level         zapcore.Level        //级别
	MaxAge        int                  //最大留存时间 天为单位
	LevelEncoder  zapcore.LevelEncoder //编码级别
	StacktraceKey string               // 栈名
	directory     string               // 输出目录
	Console       bool                 //是否在控制台打印
}

// 获取cores切片
func (z *_zap) GetZapCores() []zapcore.Core {
	cores := make([]zapcore.Core, 0, 7)
	// 将level添加到cores中
	for level := z.Level; level <= zapcore.FatalLevel; level++ {
		cores = append(cores, z.GetEncoderCore(level, z.GetLevelPriority(level)))
	}
	return cores
}

func (z *_zap) GetEncoderCore(l zapcore.Level, level zap.LevelEnablerFunc) zapcore.Core {
	writer, err := z.GetWriteSyncer(l.String()) // 使用file-rotatelogs进行日志分割
	if err != nil {
		fmt.Printf("Get Write Syncer Failed err:%v", err.Error())
		return nil
	}
	return zapcore.NewCore(z.GetEncoder(), writer, level)
}

func (z *_zap) GetWriteSyncer(level string) (zapcore.WriteSyncer, error) {
	fileWriter, err := rotatelogs.New(
		path.Join(z.directory, "%Y-%m-%d", level+".log"),
		rotatelogs.WithClock(rotatelogs.Local),
		rotatelogs.WithMaxAge(time.Duration(z.MaxAge)*24*time.Hour), // 日志留存时间
		rotatelogs.WithRotationTime(time.Hour*24),
	)
	//是否在控制台打印
	if z.Console {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(fileWriter)), err
	}
	return zapcore.AddSync(fileWriter), err
}

func (z *_zap) GetEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(z.GetEncoderConfig())
}

// 获取zapcore.EncoderConfig
func (z *_zap) GetEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  z.StacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    z.LevelEncoder,
		EncodeTime:     z.CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
}

// 自定义日志输出时间格式
func (z *_zap) CustomTimeEncoder(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	encoder.AppendString(t.Format("2006/01/02 - 15:04:05.000"))
}

// 选择调试级别
func (z *_zap) GetLevelPriority(level zapcore.Level) zap.LevelEnablerFunc {
	switch level {
	case zapcore.DebugLevel:
		return func(level zapcore.Level) bool { // 调试级别
			return level == zap.DebugLevel
		}
	case zapcore.InfoLevel:
		return func(level zapcore.Level) bool { // 日志级别
			return level == zap.InfoLevel
		}
	case zapcore.WarnLevel:
		return func(level zapcore.Level) bool { // 警告级别
			return level == zap.WarnLevel
		}
	case zapcore.ErrorLevel:
		return func(level zapcore.Level) bool { // 错误级别
			return level == zap.ErrorLevel
		}
	case zapcore.DPanicLevel:
		return func(level zapcore.Level) bool { // dpanic级别
			return level == zap.DPanicLevel
		}
	case zapcore.PanicLevel:
		return func(level zapcore.Level) bool { // panic级别
			return level == zap.PanicLevel
		}
	case zapcore.FatalLevel:
		return func(level zapcore.Level) bool { // 终止级别
			return level == zap.FatalLevel
		}
	default:
		return func(level zapcore.Level) bool { // 调试级别
			return level == zap.DebugLevel
		}
	}
}
