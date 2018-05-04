package gocongress

import (
	"context"
	"fmt"
	"os"
	"time"

	"google.golang.org/grpc"

	"github.com/go-kit/kit/log"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/spf13/cobra"
	"github.com/welaw/go-congress/client"
	"github.com/welaw/go-congress/proto"
)

func makeSendVoteCmd(grpcAddr *string) *cobra.Command {
	var limit int
	var username string
	var startDate string
	var endDate string
	var ident string

	var sendVoteCmd = &cobra.Command{
		Use:   "sendvote",
		Short: "Send new votes to welaw server.",
		Long:  "Send new votes to welaw server.",
		Run: func(cmd *cobra.Command, args []string) {
			runSendVoteCmd(*grpcAddr, ident, username, startDate, endDate, limit)
		},
	}

	sendVoteCmd.Flags().StringVar(&ident, "ident", "", "Only send vote with this identity.")
	sendVoteCmd.Flags().StringVar(&username, "username", "", "Only send votes from this user.")
	sendVoteCmd.Flags().StringVar(&startDate, "startDate", "", "Don't send votes that were published before this date.")
	sendVoteCmd.Flags().StringVar(&endDate, "endDate", "", "Don't send votes that were published after this date.")
	sendVoteCmd.Flags().IntVar(&limit, "limit", 1, "The maximum number of new votes to send.")

	return sendVoteCmd
}

func runSendVoteCmd(addr, ident, username, startDate, endDate string, limit int) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithTimeout(time.Second))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
	defer conn.Close()
	tracer := stdopentracing.GlobalTracer()
	logger := log.NewLogfmtLogger(os.Stderr)
	c := client.MakeGoCongressClient(conn, tracer, logger)
	resp, err := c.SendVote(context.TODO(), &proto.ItemRange{
		Username:  username,
		Ident:     ident,
		StartDate: startDate,
		EndDate:   endDate,
		Limit:     int32(limit),
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
	fmt.Printf("SendVote results: %v\n", len(resp.NewItems))
	for _, r := range resp.NewItems {
		fmt.Printf("%v\n", r)
	}
	return
}

func makeSendLawCmd(grpcAddr *string) *cobra.Command {
	var limit int
	var username string
	var startDate string
	var endDate string
	var ident string

	var sendLawCmd = &cobra.Command{
		Use:   "sendlaw",
		Short: "Send new laws to welaw server.",
		Long:  "Send new laws to welaw server.",
		Run: func(cmd *cobra.Command, args []string) {
			runSendLawCmd(*grpcAddr, ident, username, startDate, endDate, limit)
		},
	}

	sendLawCmd.Flags().StringVar(&ident, "ident", "", "Only send law with this identity.")
	sendLawCmd.Flags().StringVar(&username, "username", "", "Only send laws from this user.")
	sendLawCmd.Flags().StringVar(&startDate, "startDate", "", "Don't send laws that were published before this date.")
	sendLawCmd.Flags().StringVar(&endDate, "endDate", "", "Don't send laws that were published after this date.")
	sendLawCmd.Flags().IntVar(&limit, "limit", 1, "The maximum number of new laws to send.")

	return sendLawCmd
}

func runSendLawCmd(addr, ident, username, startDate, endDate string, limit int) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithTimeout(time.Second))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
	defer conn.Close()
	tracer := stdopentracing.GlobalTracer()
	logger := log.NewLogfmtLogger(os.Stderr)
	c := client.MakeGoCongressClient(conn, tracer, logger)
	resp, err := c.SendLaw(context.TODO(), &proto.ItemRange{
		Username:  username,
		Ident:     ident,
		StartDate: startDate,
		EndDate:   endDate,
		Limit:     int32(limit),
	})
	if resp != nil && resp.NewItems != nil {
		fmt.Printf("SendLaw results: %v\n", len(resp.NewItems))
		for _, r := range resp.NewItems {
			fmt.Printf("%v\n", r)
		}
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
	}
	return
}

func makeStatusCmd(grpcAddr *string) *cobra.Command {
	var username string
	var startDate string
	var endDate string
	var ident string
	var limit int

	var statusCmd = &cobra.Command{
		Use:   "status",
		Short: "get upstream status",
		Long:  `Retrieve upstream status.`,

		Run: func(cmd *cobra.Command, args []string) {
			runStatusCmd(*grpcAddr, ident, username, startDate, endDate, limit)
		},
	}

	statusCmd.Flags().StringVar(&ident, "ident", "", "Only check law with this identity.")
	statusCmd.Flags().StringVar(&username, "username", "", "Only check laws from this user.")
	statusCmd.Flags().StringVar(&startDate, "startDate", "", "Exclude laws that were published before this date.")
	statusCmd.Flags().StringVar(&endDate, "endDate", "", "Exclude laws that were published after this date.")
	statusCmd.Flags().IntVar(&limit, "limit", 1, "Stop looking if this limit is reached.")

	return statusCmd
}

func runStatusCmd(addr, ident, username, startDate, endDate string, limit int) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithTimeout(time.Second))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
	defer conn.Close()
	tracer := stdopentracing.GlobalTracer()
	logger := log.NewLogfmtLogger(os.Stderr)
	c := client.MakeGoCongressClient(conn, tracer, logger)
	reply, err := c.Status(context.TODO(), &proto.ItemRange{
		Username:  username,
		Ident:     ident,
		StartDate: startDate,
		EndDate:   endDate,
		Limit:     int32(limit),
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
	fmt.Printf("New items: %d\n", len(reply.NewItems))
	return
}

//func makeSaveCmd(grpcAddr *string) *cobra.Command {
//var format string

//var saveCmd = &cobra.Command{
//Use:     "save",
//Short:   "save files to disk",
//Long:    `Retrieve files from the upstream server.`,
//PreRunE: requiredArgs(1),
//Run: func(cmd *cobra.Command, args []string) {
//runSaveCmd(*grpcAddr, args[0], format)
//},
//}

//saveCmd.Flags().StringVar(&format, "format", "", "File format: xml, txt")

//return saveCmd
//}

//func runSaveCmd(addr, uri, filefmt string) {
//conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithTimeout(time.Second))
//if err != nil {
//fmt.Fprintf(os.Stderr, "error: %v", err)
//os.Exit(1)
//}
//defer conn.Close()
//tracer := stdopentracing.GlobalTracer()
//logger := log.NewLogfmtLogger(os.Stderr)
//c := gousa.MakeGoCongressClient(conn, tracer, logger)
//filePath, err := c.Save(context.TODO(), uri, filefmt)
//if err != nil {
//fmt.Fprintf(os.Stderr, "error: %v", err)
//os.Exit(1)
//}
//fmt.Printf("Saved: %s\n", filePath)
//return
//}
