package integration_test

import (
	ctx "context"
	"fmt"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/occam"
	"github.com/sclevine/spec"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func testSpringBoot(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		pack      occam.Pack
		docker    occam.Docker
		image     occam.Image
		container testcontainers.Container
		buildLogs fmt.Stringer
	)

	it.Before(func() {
		pack = occam.NewPack()
		docker = occam.NewDocker()
	})

	it.After(func() {
		Expect(container.Terminate(ctx.Background())).To(Succeed())
		Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
	})

	context("Maven", func() {
		it("builds SpringBoot app from source with Maven", func() {
			imageName, err := occam.RandomName()
			Expect(err).ToNot(HaveOccurred())

			image, buildLogs, err = pack.WithNoColor().Build.
				WithBuildpacks(buildPack).
				WithEnv(map[string]string{
					"BP_ARCH": "amd64",
				}).
				WithBuilder(builder).
				WithTrustBuilder().
				WithPullPolicy("if-not-present").
				Execute(imageName, "samples/java/maven")
			Expect(err).ToNot(HaveOccurred())
			Expect(buildLogs.String()).ToNot(BeEmpty())
			Expect(len(image.Buildpacks)).To(BeNumerically(">", 0))

			container, err = testContainers.WithExposedPorts("8080/tcp").WithWaitingFor(wait.ForLog("Started DemoApplication in")).Execute(imageName)
			Expect(err).NotTo(HaveOccurred())
			mappedPort, err := container.MappedPort(ctx.Background(), "8080/tcp")
			Expect(err).ShouldNot(HaveOccurred())
			resp, err := makeRequest("/actuator/health", mappedPort.Port())
			Expect(err).ShouldNot(HaveOccurred())
			defer resp.Body.Close()
			Expect(resp.StatusCode).To(Equal(200))
		})
	})

	context("Gradle", func() {
		it("builds SpringBoot app from source with Gradle", func() {
			imageName, err := occam.RandomName()
			Expect(err).ToNot(HaveOccurred())

			image, buildLogs, err = pack.WithNoColor().Build.
				WithBuildpacks(buildPack).
				WithEnv(map[string]string{
					"BP_ARCH": "amd64",
				}).
				WithBuilder(builder).
				WithTrustBuilder().
				WithPullPolicy("if-not-present").
				Execute(imageName, "samples/java/gradle")
			Expect(err).ToNot(HaveOccurred())
			Expect(buildLogs.String()).ToNot(BeEmpty())
			Expect(len(image.Buildpacks)).To(BeNumerically(">", 0))

			container, err = testContainers.WithExposedPorts("8080/tcp").WithWaitingFor(wait.ForLog("Started DemoApplication in")).Execute(imageName)
			Expect(err).NotTo(HaveOccurred())
			mappedPort, err := container.MappedPort(ctx.Background(), "8080/tcp")
			Expect(err).ShouldNot(HaveOccurred())
			resp, err := makeRequest("/actuator/health", mappedPort.Port())
			Expect(err).ShouldNot(HaveOccurred())
			defer resp.Body.Close()
			Expect(resp.StatusCode).To(Equal(200))
		})
	})
}
