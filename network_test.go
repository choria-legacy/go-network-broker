package network

import (
	"os"
	"testing"
	"time"

	"github.com/choria-io/go-choria/build"
	"github.com/choria-io/go-choria/choria"

	"github.com/choria-io/go-choria/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func Test(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Choria Network Broker")
}

var _ = Describe("Network Broker", func() {
	var (
		cfg *config.Config
		fw  *choria.Framework
		srv *Server
		err error
	)

	BeforeEach(func() {
		os.Setenv("MCOLLECTIVE_CERTNAME", "rip.mcollective")

		cfg, err = config.NewDefaultConfig()
		Expect(err).ToNot(HaveOccurred())

		cfg.Choria.SSLDir = "testdata/ssl"

		fw, err = choria.NewWithConfig(cfg)
		Expect(err).ToNot(HaveOccurred())
	})

	var _ = Describe("NewServer", func() {
		It("Should initialize the server correctly", func() {
			cfg.Choria.NetworkListenAddress = "0.0.0.0"
			cfg.Choria.NetworkClientPort = 8080
			cfg.Choria.NetworkWriteDeadline = time.Duration(10 * time.Second)
			cfg.LogLevel = "error"
			cfg.Choria.StatsPort = 8081
			cfg.Choria.StatsListenAddress = "192.168.1.1"
			cfg.Choria.NetworkPeerPort = 8082
			cfg.Choria.NetworkPeerUser = "bob"
			cfg.Choria.NetworkPeerPassword = "secret"
			cfg.Choria.NetworkPeers = []string{"nats://localhost:9000", "nats://localhost:9001", "nats://localhost:8082"}

			fw, err = choria.NewWithConfig(cfg)
			Expect(err).ToNot(HaveOccurred())

			srv, err = NewServer(fw, false)
			Expect(err).ToNot(HaveOccurred())

			Expect(srv.opts.Host).To(Equal("0.0.0.0"))
			Expect(srv.opts.Port).To(Equal(8080))
			Expect(srv.opts.Logtime).To(BeFalse())
			Expect(srv.opts.MaxConn).To(Equal(build.MaxBrokerClients()))
			Expect(srv.opts.WriteDeadline).To(Equal(time.Duration(10 * time.Second)))
			Expect(srv.opts.NoSigs).To(BeTrue())
			Expect(srv.opts.Debug).To(BeFalse())
			Expect(srv.opts.HTTPHost).To(Equal("192.168.1.1"))
			Expect(srv.opts.HTTPPort).To(Equal(8081))
			Expect(srv.opts.Cluster.Host).To(Equal("0.0.0.0"))
			Expect(srv.opts.Cluster.NoAdvertise).To(BeTrue())
			Expect(srv.opts.Cluster.Port).To(Equal(8082))
			Expect(srv.opts.Cluster.Username).To(Equal("bob"))
			Expect(srv.opts.Cluster.Password).To(Equal("secret"))
			Expect(srv.opts.Routes).To(HaveLen(2))
			Expect(srv.opts.Routes[0].Host).To(Equal("localhost:9000"))
			Expect(srv.opts.Routes[1].Host).To(Equal("localhost:9001"))
			Expect(srv.opts.TLS).To(BeTrue())
			Expect(srv.opts.TLSVerify).To(BeTrue())
			Expect(srv.opts.TLSTimeout).To(Equal(float64(2)))
			Expect(srv.opts.Cluster.TLSTimeout).To(Equal(float64(2)))
		})

		It("Should support disabling TLS Verify", func() {
			cfg.DisableTLSVerify = true

			fw, err = choria.NewWithConfig(cfg)
			Expect(err).ToNot(HaveOccurred())

			srv, err = NewServer(fw, false)
			Expect(err).ToNot(HaveOccurred())
			Expect(srv.opts.TLSVerify).To(BeFalse())
		})

		It("Should support disabling TLS", func() {
			cfg.DisableTLS = true

			fw, err = choria.NewWithConfig(cfg)
			Expect(err).ToNot(HaveOccurred())

			srv, err = NewServer(fw, false)
			Expect(err).ToNot(HaveOccurred())
			Expect(srv.opts.TLS).To(BeFalse())
		})
	})
})