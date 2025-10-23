package gather

import (
	"context"
	"encoding/hex"
	"errors"
	"net"
	"time"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	ctx               context.Context
	servers           []string
	conn              *grpc.ClientConn
	timeout           time.Duration
	blockchainClient  pactus.BlockchainClient
	transactionClient pactus.TransactionClient
}

func NewGrpcClient(timeout time.Duration, servers []string) *GrpcClient {
	ctx := context.Background()

	cli := &GrpcClient{
		ctx:               ctx,
		timeout:           timeout,
		conn:              nil,
		blockchainClient:  nil,
		transactionClient: nil,
	}

	if len(servers) > 0 {
		cli.servers = servers
	}

	return cli
}

func (c *GrpcClient) Connect() error {
	if c.conn != nil {
		return nil
	}

	for _, server := range c.servers {
		conn, err := grpc.NewClient(server,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithContextDialer(func(_ context.Context, s string) (net.Conn, error) {
				return net.DialTimeout("tcp", s, c.timeout)
			}))
		if err != nil {
			continue
		}

		blockchainClient := pactus.NewBlockchainClient(conn)
		transactionClient := pactus.NewTransactionClient(conn)

		// Check if client is responding
		_, err = blockchainClient.GetBlockchainInfo(c.ctx,
			&pactus.GetBlockchainInfoRequest{})
		if err != nil {
			_ = conn.Close()

			continue
		}

		c.conn = conn
		c.blockchainClient = blockchainClient
		c.transactionClient = transactionClient

		return nil
	}

	return errors.New("unable to connect to the servers")
}

func (c *GrpcClient) GetBlockchainInfo() (*pactus.GetBlockchainInfoResponse, error) {
	if err := c.Connect(); err != nil {
		return nil, err
	}

	info, err := c.blockchainClient.GetBlockchainInfo(c.ctx,
		&pactus.GetBlockchainInfoRequest{})
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (c *GrpcClient) GetBlock(height uint32, verbosity pactus.BlockVerbosity) (*pactus.GetBlockResponse, error) {
	if err := c.Connect(); err != nil {
		return nil, err
	}

	info, err := c.blockchainClient.GetBlock(c.ctx,
		&pactus.GetBlockRequest{Height: height, Verbosity: verbosity})
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (c *GrpcClient) getAccount(addrStr string) (*pactus.AccountInfo, error) {
	if err := c.Connect(); err != nil {
		return nil, err
	}

	res, err := c.blockchainClient.GetAccount(c.ctx,
		&pactus.GetAccountRequest{Address: addrStr})
	if err != nil {
		return nil, err
	}

	return res.Account, nil
}

func (c *GrpcClient) getValidator(addrStr string) (*pactus.ValidatorInfo, error) {
	if err := c.Connect(); err != nil {
		return nil, err
	}

	res, err := c.blockchainClient.GetValidator(c.ctx,
		&pactus.GetValidatorRequest{Address: addrStr})
	if err != nil {
		return nil, err
	}

	return res.Validator, nil
}

func (c *GrpcClient) sendTx(trx *tx.Tx) (tx.ID, error) {
	if err := c.Connect(); err != nil {
		return hash.UndefHash, err
	}

	data, err := trx.Bytes()
	if err != nil {
		return hash.UndefHash, err
	}
	res, err := c.transactionClient.BroadcastTransaction(c.ctx,
		&pactus.BroadcastTransactionRequest{SignedRawTransaction: hex.EncodeToString(data)})
	if err != nil {
		return hash.UndefHash, err
	}

	return hash.FromString(res.Id)
}

// TODO: check the return value type.
func (c *GrpcClient) getTransaction(id tx.ID) (*pactus.GetTransactionResponse, error) {
	if err := c.Connect(); err != nil {
		return nil, err
	}

	res, err := c.transactionClient.GetTransaction(c.ctx,
		&pactus.GetTransactionRequest{
			Id:        id.String(),
			Verbosity: pactus.TransactionVerbosity_TRANSACTION_VERBOSITY_INFO,
		})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *GrpcClient) getFee(amt amount.Amount, payloadType payload.Type) (amount.Amount, error) {
	if err := c.Connect(); err != nil {
		return 0, err
	}

	res, err := c.transactionClient.CalculateFee(c.ctx,
		&pactus.CalculateFeeRequest{
			Amount:      amt.ToNanoPAC(),
			PayloadType: pactus.PayloadType(payloadType),
		})
	if err != nil {
		return 0, err
	}

	return amount.Amount(res.Fee), nil
}
