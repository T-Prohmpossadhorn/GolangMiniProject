package internal

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/T-Prohmpossadhorn/GolangMiniProject/internal/config"
	"github.com/T-Prohmpossadhorn/GolangMiniProject/internal/db"
	"github.com/T-Prohmpossadhorn/GolangMiniProject/internal/go-grpcserver"
	"github.com/T-Prohmpossadhorn/GolangMiniProject/internal/middleware"
	"github.com/T-Prohmpossadhorn/GolangMiniProject/internal/persist"
	"github.com/T-Prohmpossadhorn/GolangMiniProject/internal/transactionservice"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/encoding/protojson"

	proto "github.com/T-Prohmpossadhorn/GolangMiniProject/pkg/Proto"
)

var timeout = 10 * time.Second

type Service struct {
	srv *grpcserver.Server
	db  persist.Persist
}

func New(ctx context.Context, options config.Options) (*Service, error) {
	listener, err := net.Listen("tcp", options.ListenAddress)
	if err != nil {
		return nil, fmt.Errorf("net.Listen(): %w", err)
	}

	svc := Service{}

	DB := db.New(options.FilePath)

	svc.db = DB
	svc.srv = grpcserver.New()

	grpcSrv, err := svc.srv.Create()
	if err != nil {
		log.Fatal("svc.srv.Create(): %w", err)
	}

	transactionsrv, err := transactionservice.New(DB, options)
	if err != nil {
		return nil, fmt.Errorf("transactionsrv.New(): %w", err)
	}

	proto.RegisterFruitListServiceServer(grpcSrv, transactionsrv)

	go func() {
		if err := grpcSrv.Serve(listener); err != nil {
			log.Fatal("grpcSrv.Serve(): %w", err)
		}
	}()

	if options.ListenAddressHTTP != "" {
		if err := CreateHTTPServer(ctx, transactionsrv, options); err != nil {
			log.Fatal("CreateHTTPServer(): %w", err)
		}
	}

	return &svc, nil
}

func CreateHTTPServer(ctx context.Context, transactionsrv *transactionservice.TransactionService, options config.Options) error {
	mux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			Multiline:       true,
			Indent:          "",
			AllowPartial:    false,
			UseProtoNames:   false,
			UseEnumNumbers:  false,
			EmitUnpopulated: true,
			Resolver:        nil,
		},
	}))

	if err := proto.RegisterFruitListServiceHandlerServer(ctx, mux, transactionsrv); err != nil {
		return fmt.Errorf("proto.RegisterFinancialTransactionServiceHandler(): %w", err)
	}

	h := middleware.JSONContentTypeValidator(mux)

	server := http.Server{
		Addr:           options.ListenAddressHTTP,
		Handler:        h,
		ReadTimeout:    timeout,
		WriteTimeout:   timeout,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
	}

	if err := server.ListenAndServe(); err != nil {
		return fmt.Errorf("http.ListenAndServe(): %w", err)
	}

	return nil
}

func (s *Service) Shutdown() bool {
	return s.srv.Shutdown()
}
