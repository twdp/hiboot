// Copyright 2018 John Deng (hi.devops.io@gmail.com).
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package autoconfigure_test

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/factory"
	"github.com/hidevopsio/hiboot/pkg/factory/autoconfigure"
	"github.com/hidevopsio/hiboot/pkg/factory/instantiate"
	"github.com/hidevopsio/hiboot/pkg/inject"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils/cmap"
	"github.com/hidevopsio/hiboot/pkg/utils/io"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

type FakeProperties struct {
	Name     string `default:"foo"`
	Nickname string `default:"foobar"`
	Username string `default:"fb"`
	Org      string `default:"hidevopsio"`
	Profile  string `default:"${APP_PROFILES_ACTIVE:dev}"`
}

type unknownConfiguration struct {
	FakeProperties FakeProperties `mapstructure:"fake"`
}

type FooInterface interface {
}

type unsupportedConfiguration struct {
	FooInterface
	FakeProperties FakeProperties `mapstructure:"fake"`
}

type FooProperties struct {
	Name     string `default:"${fake.name}"`
	Nickname string `default:"foobar"`
	Username string `default:"fb"`
}

type Configuration struct {
	at.AutoConfiguration
}

type emptyConfiguration struct {
}

type FooConfiguration struct {
	at.AutoConfiguration
	FakeProperties FooProperties `mapstructure:"foo"`
}

type BarConfiguration struct {
	app.Configuration
	FakeProperties FooProperties `mapstructure:"bar"`
}

type FooBarConfiguration struct {
	app.Configuration
	FakeProperties FooProperties `mapstructure:"foobar"`
	foobar         *FooBar
}

type mercury struct {
	jupiter *jupiter
}

type mars struct {
	mercury *mercury
}

type jupiter struct {
}

type marsConfiguration struct {
	app.Configuration `depends:"mercuryConfiguration"`
}

func (c *marsConfiguration) Mars(m *mercury) *mars {
	return &mars{mercury: m}
}

type mercuryConfiguration struct {
	app.Configuration `depends:"jupiterConfiguration"`
}

func (c *mercuryConfiguration) Mercury(j *jupiter) *mercury {
	return &mercury{jupiter: j}
}

type jupiterConfiguration struct {
	app.Configuration
}

func (c *jupiterConfiguration) Jupiter() *jupiter {
	return new(jupiter)
}

func newFooBarConfiguration(foobar *FooBar) *FooBarConfiguration {
	return &FooBarConfiguration{
		foobar: foobar,
	}
}

type HelloWorld string
type Hello string

func (c *FooConfiguration) HelloWorld(foo Hello) HelloWorld {
	return HelloWorld(foo + " world")
}

func (c *FooConfiguration) Hello() Hello {
	return Hello("Hello")
}

type Foo struct {
	Name string
	Bar  *Bar
}

type Bar struct {
	Name string
}

type FooBar struct {
	Name string
}

type barConfiguration struct {
	FakeProperties FooProperties `mapstructure:"bar"`
}

func (c *barConfiguration) BarBar() *Bar {
	return new(Bar)
}

func init() {
	log.SetLevel(log.DebugLevel)
	io.EnsureWorkDir(1, "config/application.yml")
}

type fakeConfiguration struct {
	app.Configuration
	FakeProperties FakeProperties `mapstructure:"fake"`
}

type foobarConfiguration struct {
	app.Configuration
	FakeProperties FakeProperties `mapstructure:"foobar"`
}

func newFoobarConfiguration() *foobarConfiguration {
	return &foobarConfiguration{}
}

func (c *foobarConfiguration) Foo(bar *Bar) *Foo {
	f := new(Foo)
	f.Name = c.FakeProperties.Name
	f.Bar = bar
	return f
}

func (c *foobarConfiguration) FooBar(foo *Foo) *FooBar {
	f := new(FooBar)
	f.Name = foo.Name

	return f
}

func (c *foobarConfiguration) Bar() *Bar {
	b := new(Bar)
	b.Name = "bar"

	return b
}

type EarthConfiguration struct {
	app.Configuration
	instantiateFactory factory.InstantiateFactory
}

func newEarthConfiguration(instantiateFactory factory.InstantiateFactory) *EarthConfiguration {
	c := &EarthConfiguration{
		instantiateFactory: instantiateFactory,
	}
	c.RuntimeDeps.Set(c.RuntimeTree, []string{"autoconfigure_test.leaf", "autoconfigure_test.branch"})
	return c
}

type Land struct {
	Mountain *Mountain
}

type RuntimeTree struct {
	leaf   *Leaf
	branch *Branch
}

type Tree interface {
	Family() string
}

type oakTree struct {
	Tree
}

func (t *oakTree) Family() string {
	return "Beech"
}

type mapleTree struct {
}

func (t *mapleTree) Family() string {
	return "Soapberry"
}

type Branch struct {
}

type Leaf struct {
}

type Mountain struct {
	oakTree Tree
}

func (c *EarthConfiguration) Land(mountain *Mountain) *Land {
	return &Land{Mountain: mountain}
}

func (c *EarthConfiguration) Mountain(tree Tree) *Mountain {
	return &Mountain{oakTree: tree}
}

func (c *EarthConfiguration) Branch() *Branch {
	return &Branch{}
}

func (c *EarthConfiguration) Leaf() *Leaf {
	return &Leaf{}
}

func (c *EarthConfiguration) RuntimeTree() *RuntimeTree {
	leaf := c.instantiateFactory.GetInstance("autoconfigure_test.leaf").(*Leaf)
	branch := c.instantiateFactory.GetInstance("autoconfigure_test.branch").(*Branch)
	return &RuntimeTree{leaf: leaf, branch: branch}
}

func (c *EarthConfiguration) OakTree() Tree {
	return &oakTree{}
}

func (c *EarthConfiguration) MapleTree() Tree {
	return &mapleTree{}
}

type helloService struct{ foo *Foo }

func newHelloService(foo *Foo) *helloService {
	return &helloService{foo: foo}
}

func TestConfigurableFactory(t *testing.T) {

	configPath := filepath.Join(os.TempDir(), "config")

	fakeFile := "application.yml"
	os.Remove(filepath.Join(configPath, fakeFile))
	fakeContent :=
		"app:\n" +
			"  project: hidevopsio\n" +
			"  name: hiboot\n" +
			"  version: ${unknown.version:0.0.1}\n"
	n, err := io.WriterFile(configPath, fakeFile, []byte(fakeContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, n, len(fakeContent))

	configContainers := cmap.New()
	customProperties := cmap.New()

	f := autoconfigure.NewConfigurableFactory(
		instantiate.NewInstantiateFactory(cmap.New(), []*factory.MetaData{}, customProperties),
		configContainers,
	)

	t.Run("should set factory instance", func(t *testing.T) {
		err = f.SetInstance(factory.InstantiateFactoryName, f)
		assert.Equal(t, nil, err)

		f.SetInstance(factory.ConfigurableFactoryName, f)
		assert.Equal(t, nil, err)
	})
	inject.SetFactory(f)

	io.ChangeWorkDir(os.TempDir())
	t.Run("should build app config", func(t *testing.T) {
		// backup profile
		profile := os.Getenv(autoconfigure.EnvAppProfilesActive)
		os.Setenv(autoconfigure.EnvAppProfilesActive, "")
		sc, err := f.BuildSystemConfig()
		assert.Equal(t, nil, err)
		assert.Equal(t, "default", sc.App.Profiles.Active)
		// restore profile
		os.Setenv(autoconfigure.EnvAppProfilesActive, profile)
	})

	t.Run("should build app config", func(t *testing.T) {
		// backup profile
		customProperties.Set("app.profiles.active", "dev")
		customProperties.Set("logging.level", "debug")
		profile := os.Getenv(autoconfigure.EnvAppProfilesActive)
		os.Setenv(autoconfigure.EnvAppProfilesActive, "")
		sc, err := f.BuildSystemConfig()
		assert.Equal(t, nil, err)
		assert.Equal(t, "dev", sc.App.Profiles.Active)
		// restore profile
		os.Setenv(autoconfigure.EnvAppProfilesActive, profile)
	})

	t.Run("should build app config", func(t *testing.T) {
		_, err = f.BuildSystemConfig()
		assert.Equal(t, nil, err)
	})

	fooFile := "application-foo.yml"
	os.Remove(filepath.Join(configPath, fooFile))
	fooContent :=
		"foo:\n" +
			"  name: foo\n" +
			"  nickname: ${app.name} ${foo.name}\n" +
			"  username: ${unknown.name:bar}\n"
	_, err = io.WriterFile(configPath, fooFile, []byte(fooContent))
	assert.Equal(t, nil, err)

	fooConfig := new(FooConfiguration)
	fakeCfg := new(fakeConfiguration)

	f.Build([]*factory.MetaData{
		factory.NewMetaData(new(emptyConfiguration)),
		factory.NewMetaData("foo", fooConfig),
		factory.NewMetaData(fakeCfg),
		factory.NewMetaData(new(Configuration)),
		factory.NewMetaData(new(BarConfiguration)),
		factory.NewMetaData(new(marsConfiguration)),
		factory.NewMetaData(new(jupiterConfiguration)),
		factory.NewMetaData(new(mercuryConfiguration)),
		factory.NewMetaData(new(unknownConfiguration)),
		factory.NewMetaData(new(unsupportedConfiguration)),
		factory.NewMetaData(newFoobarConfiguration),
		factory.NewMetaData(newEarthConfiguration),
	})

	f.AppendComponent(newHelloService)

	f.BuildComponents()

	t.Run("should get SystemConfiguration", func(t *testing.T) {
		sysCfg := f.SystemConfiguration()
		assert.NotEqual(t, nil, sysCfg)
	})

	t.Run("should get fake configuration", func(t *testing.T) {
		fakeCfg := f.Configuration("fake")
		assert.NotEqual(t, nil, fakeCfg)
	})

	t.Run("should not get non-existence configuration", func(t *testing.T) {
		fakeCfg := f.Configuration("non-existence")
		assert.Equal(t, nil, fakeCfg)
	})

	t.Run("should add instance to factory at runtime", func(t *testing.T) {
		fakeInstance := &struct{ Name string }{Name: "fake"}
		f.SetInstance("autoconfigure_test.fakeInstance", fakeInstance)
		gotFakeInstance := f.GetInstance("autoconfigure_test.fakeInstance")
		assert.Equal(t, fakeInstance, gotFakeInstance)
	})

	t.Run("should get foo configuration", func(t *testing.T) {
		helloWorld := f.GetInstance("autoconfigure_test.helloWorld")
		assert.NotEqual(t, nil, helloWorld)
		assert.Equal(t, HelloWorld("Hello world"), helloWorld)

		assert.Equal(t, "hiboot foo", fooConfig.FakeProperties.Nickname)
		assert.Equal(t, "bar", fooConfig.FakeProperties.Username)
		assert.Equal(t, "foo", fooConfig.FakeProperties.Name)
	})

	t.Run("should get runtime created instances", func(t *testing.T) {
		runtimeTree := f.GetInstance(RuntimeTree{})
		leaf := f.GetInstance(Leaf{})
		branch := f.GetInstance(Branch{})
		assert.NotEqual(t, nil, runtimeTree)
		rt := runtimeTree.(*RuntimeTree)
		assert.Equal(t, branch, rt.branch)
		assert.Equal(t, leaf, rt.leaf)
	})
}
