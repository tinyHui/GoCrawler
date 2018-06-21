# GoLang based web crawler

## How to Run
1. Install dependencies `make install`
2. Build `make build`
3. Run `MODE=<enviroment> config=<config file path> <inital URL>`
4. Wait it finish and open the sitemap file, done

## Requirements
1. go ~1.9.3
2. glide ~0.13.1
3. make 3.81

## Configuration
- A configuration file required during executing, the format can be found inside `./config/prod/parameters.yaml`.
- Executing environment is required, need to be set as an environment variable, the variable name is `MODE`. Default environment is taken as `debug`.
- You need to define the sitemap save path on your machine.