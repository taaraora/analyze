package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/supergiant/analyze/pkg/logger"

	"github.com/aws/aws-sdk-go-v2/aws/ec2metadata"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func main() {
	var (
		port      = flag.Int("api-port", 9292, "tcp port where node agent API is serving")
		logLevel  = flag.String("log-level", "debug", "logging level, e.g. info, warning, debug, error, fatal")
		logFormat = flag.String("log-format", "TXT", "logging format [TXT JSON]")
	)

	flag.Parse()

	loggerConf := logger.Config{
		Level:     *logLevel,
		Formatter: logger.Formatter(*logFormat),
	}

	if err := loggerConf.Validate(); err != nil {
		log.Fatalf("\n%v\n", err)
	}

	if *port < 1 {
		log.Fatalf("api-port is %v, it need to be greater than 0", *port)
	}

	logger := logger.NewLogger(loggerConf).WithField("app", "analyze-nodeagent")

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		logger.Fatalf("unable to load SDK config, %v", err)
	}

	var ec2MetadataService = ec2metadata.New(cfg)

	if strings.ToLower(os.Getenv("AWS_EC2_METADATA_DISABLED")) == "true" {
		logger.Error("AWS_EC2_METADATA_DISABLED is true, we are not able to interact with EC2 metadata API")
	}

	if available := ec2MetadataService.Available(); !available {
		logger.Error("metadata is NOT available")
	}

	result, err := ec2MetadataService.GetMetadata("instance-id")
	if err != nil {
		logger.Error(err)
	}

	logger.Infof("metadata is available, at instance-id: %s", result)

	var httpServer = &http.Server{}
	var router = mux.NewRouter()
	awsAPI := router.PathPrefix("/aws").Subrouter()
	httpServer.Handler = awsAPI

	awsAPI.HandleFunc("/meta-data/{path}", func(
		ec2MetadataService *ec2metadata.EC2Metadata,
		logger logrus.FieldLogger,
	) func(http.ResponseWriter, *http.Request) {
		return func(res http.ResponseWriter, req *http.Request) {
			vars := mux.Vars(req)

			path := vars["path"]
			result, err := ec2MetadataService.GetMetadata(path)
			if err != nil {
				logger.Errorf("can't make GetMetadata request to ec2 metadata api, %+v", err)
			}
			_, err = res.Write([]byte(result))
			if err != nil {
				logger.Errorf("can't write ec2 metadata content, %+v", err)
			}
		}
	}(ec2MetadataService, logger)).Methods(http.MethodGet)

	awsAPI.HandleFunc("/dynamic/{path}", func(
		ec2MetadataService *ec2metadata.EC2Metadata,
		logger logrus.FieldLogger,
	) func(http.ResponseWriter, *http.Request) {
		return func(res http.ResponseWriter, req *http.Request) {
			vars := mux.Vars(req)

			path := vars["path"]
			result, err := ec2MetadataService.GetDynamicData(path)
			if err != nil {
				logger.Errorf("can't make GetDynamicData request to ec2 metadata api, %+v", err)
			}
			_, err = res.Write([]byte(result))
			if err != nil {
				logger.Errorf("can't write ec2 metadata content, $+v", err)
			}
		}
	}(ec2MetadataService, logger)).Methods(http.MethodGet)

	awsAPI.HandleFunc("/user-data", func(
		ec2MetadataService *ec2metadata.EC2Metadata,
		logger logrus.FieldLogger,
	) func(http.ResponseWriter, *http.Request) {
		return func(res http.ResponseWriter, req *http.Request) {
			result, err := ec2MetadataService.GetUserData()
			if err != nil {
				logger.Errorf("can't make GetUserData request to ec2 metadata api %+v", err)
			}
			_, err = res.Write([]byte(result))
			if err != nil {
				logger.Errorf("can't write ec2 metadata content, %+v", err)
			}
		}
	}(ec2MetadataService, logger)).Methods(http.MethodGet)

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(*port))
	if err != nil {
		logger.Fatalf("unable to start api listener: %v", err)
	}

	var addr = listener.Addr().String()
	logger.Infof("listener has started, address: %s", addr)

	if err := httpServer.Serve(listener); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("unable to start api listener: %v", err)
	}

}
