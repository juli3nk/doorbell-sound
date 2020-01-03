package main

type Options struct {
	CORSAllowOrigins []string `split_words:"true" default:"*"`
	CORSAllowMethods []string `split_words:"true" default:"GET"`
}
