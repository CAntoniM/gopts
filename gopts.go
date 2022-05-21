package gopts

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

const (
	// ANSII terminal format codes used to make the terminal output look pretty
	Red    = "\033[31m"
	Bold   = "\033[1m"
	Reset  = "\033[0m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
	Green  = "\033[32m"
	// Enumeration of the types of elements inside of api
	sflag_element    = 0
	lflag_element    = 1
	flag_element     = 2
	required_element = 3
	optional_element = 4
)

// Internal data used to determin every thing needed to parse the arguments
type flag_info struct {
	long        bool
	short       bool
	takes_value bool
	desc        string
}

// struct to encapsulate all of the required infomation to parse the optional
// and required arguments
type value_info struct {
	name string
	desc string
}

// This is the base object that used to hold all of the infomation that is used
// to parse the arguments. Users are expected to get this via a call to the Opts
// function and are expected to leave the internals alone and to instread use
// the methods provided to build up the options parsing
type Options struct {
	flags           map[string]flag_info
	optional_values []value_info
	required_values []value_info
	description     string
}

// This the constructor function for the arguments parser. this method takes in
// a string that will be used as the description text as part of the help text.
func Opts(description string) *Options {
	opt := Options{}
	opt.description = description
	opt.flags = make(map[string]flag_info)
	return &opt
}

// This defines to the application to add an optional argument to the command
// line interface such that the value that will places into the flags map that
// is returned form the parse call is either empty, "true", or the value that
// was given.
func (opt *Options) Optional(name string, desc string) *Options {

	if name == "" {
		panic(errors.New("please ensure that all flags contain a name"))
	}

	for _, info := range opt.optional_values {
		if info.name == name {
			panic(fmt.Errorf("invalid optional value added 2 items can not share the name %s", name))
		}
	}
	opt.optional_values = append(opt.optional_values, value_info{name, desc})

	return opt
}

//This adds a commandline argument to the options build that must be present for
//the application to start if this element is not present then the application
// print an error message and exit.
func (opt *Options) Required(name string, desc string) *Options {

	if name == "" {
		panic(errors.New("please ensure that all values contain a name"))
	}

	for _, info := range opt.required_values {
		if info.name == name {
			panic(fmt.Errorf("invalid required value added 2 items can not share the name %s", name))
		}
	}
	opt.required_values = append(opt.required_values, value_info{name, desc})

	return opt
}

//This defines a flag that takes both a short and long form of the flag e.g. -a,--addr
//this will add this element to the builder to be used in the parse call
func (opt *Options) Flag(name string, desc string, value ...bool) *Options {

	var _value bool = true
	for _, value := range value {
		_value = value
	}

	if name == "" {
		panic(errors.New("please ensure that all flags contain a name"))
	}

	if opt.flags[name] != (flag_info{}) {
		panic(fmt.Errorf("invalid flag added 2 flags can not share the name %s", name))
	}
	opt.flags[name] = flag_info{true, true, _value, desc}

	return opt
}

//This defines a flag that takes both a short form of the flag e.g. -a
//this will add this element to the builder to be used in the parse call
func (opt *Options) Short_flag(name string, desc string, value ...bool) *Options {

	var _value bool = true
	for _, value := range value {
		_value = value
	}

	if name == "" {
		panic(errors.New("please ensure that all flags contain a name"))
	}

	for key, value := range opt.flags {
		if key[0] == name[0] && value.short {
			panic(fmt.Errorf("invalid the short flag %c has already been defined please define another flag for this value", name[0]))
		}
	}

	info := flag_info{false, true, _value, desc}
	if opt.flags[name] != (flag_info{}) {
		info = opt.flags[name]
		info.short = true
	}
	opt.flags[name] = info

	return opt
}

//This defines a flag that takes both a long form of the flag e.g. --addr
//this will add this element to the builder to be used in the parse call
func (opt *Options) Long_flag(name string, desc string, value ...bool) *Options {
	var _value bool = true
	for _, value := range value {
		_value = value
	}
	if name == "" {
		panic(errors.New("please ensure that all flags contain a name"))
	}

	if opt.flags[name] == (flag_info{}) {
		panic(fmt.Errorf("invalid flag added 2 flags can not share the name %s", name))
	}

	info := flag_info{true, false, _value, desc}
	if opt.flags[name] != (flag_info{}) {
		info = opt.flags[name]
		info.long = true
	}
	opt.flags[name] = info

	return opt
}

//This defines that a error has accured during the parsing of the arguments and
// that we need to tell the user in a pretty and consistent way we should if possible also given
// the user a hint on how to fix it. this will exit the application
func (opt *Options) Err(message string, hint string) {

	os.Stderr.Write([]byte(fmt.Sprintf("%s%sERROR%s: %s\n\n", Red, Bold, Reset, message)))

	if hint != "" {
		os.Stderr.Write([]byte(fmt.Sprintf("\t%s%sHint%s: %s\n\n", Yellow, Bold, Reset, hint)))
	}

	os.Exit(-1)
}

func (opt *Options) Help() {
	procname := filepath.Base(os.Args[0])
	os.Stdout.Write([]byte(fmt.Sprintf("%s%s%s%s \n\t%s\n\n", Blue, Bold, procname, Reset, opt.description)))
	os.Stdout.Write([]byte(fmt.Sprintf("%s%sUSAGE:%s\n\t%s ", Cyan, Bold, Reset, procname)))

	if len(opt.flags) != 0 {
		os.Stdout.Write([]byte("[OPTIONS] "))
	}

	for _, info := range opt.required_values {
		os.Stdout.Write([]byte(fmt.Sprintf("<%s> ", strings.ToUpper(info.name))))
	}
	for _, info := range opt.optional_values {
		os.Stdout.Write([]byte(fmt.Sprintf("[%s] ", strings.ToUpper(info.name))))
	}

	os.Stdout.Write([]byte("\n\n"))
	if len(opt.optional_values) != 0 || len(opt.required_values) != 0 {
		os.Stdout.Write([]byte(fmt.Sprintf("%s%sARGS:%s\n", Cyan, Bold, Reset)))
		var max_len int
		for _, info := range opt.required_values {
			if len(info.name) > max_len {
				max_len = len(info.name)
			}
		}

		for _, info := range opt.optional_values {
			if len(info.name) > max_len {
				max_len = len(info.name)
			}
		}

		for _, info := range opt.required_values {
			diff := max_len - len(info.name)
			i := 0
			spacer := ""
			for i < diff {
				spacer = spacer + " "
			}
			os.Stdout.Write([]byte(fmt.Sprintf("\t%s<%s>%s    %s%s\n", Blue, strings.ToUpper(info.name), Reset, spacer, info.desc)))
		}

		for _, info := range opt.optional_values {
			diff := max_len - len(info.name)
			spacer := ""
			for i := 0; i < diff; i++ {
				spacer = spacer + " "
			}
			os.Stdout.Write([]byte(fmt.Sprintf("\t%s<%s>%s    %s%s\n", Blue, strings.ToUpper(info.name), Reset, spacer, info.desc)))
		}
	}

	if len(opt.flags) != 0 {
		os.Stdout.Write([]byte(fmt.Sprintf("%s%sOPTIONS:%s\n", Cyan, Bold, Reset)))
		var max_len int
		for key := range opt.flags {
			if len(key) > max_len {
				max_len = len(key)
			}
		}
		for key, value := range opt.flags {
			os.Stdout.Write([]byte(fmt.Sprintf("\t%s", Blue)))
			if value.short {
				os.Stdout.Write([]byte(fmt.Sprintf("-%c,", key[0])))
			} else {
				os.Stdout.Write([]byte("   "))
			}
			diff := max_len - len(key)
			spacer := ""
			if value.long {
				for i := 0; i < diff; i++ {
					spacer = spacer + " "
				}
				os.Stdout.Write([]byte(fmt.Sprintf(" --%s%s %s    %s\n", key, Reset, spacer, value.desc)))
			} else {
				for i := 0; i < max_len; i++ {
					spacer = spacer + " "
				}
				os.Stdout.Write([]byte(fmt.Sprintf("   %s %s    %s\n", Reset, spacer, value.desc)))
			}
		}
	}
	os.Exit(0)
}

//helper function that returns the name and flag info of a flag that meets to criteria
// to match the short flag given
func (opt *Options) Get_sflag(name rune) (string, flag_info) {
	for key, value := range opt.flags {
		if rune(key[0]) == name && value.short {
			return key, value
		}
	}
	return "", flag_info{}
}

// This function will take the builder provided and transform it into three maps
// these are:
// 		required: this defines all arguments added via the required function
//		optional: this defines all arguments added via the optional function
//		flags: this defines all flags give on the execution of the application
//
// It will do this vai using the given the infomation inside of the builder
// and it will use it structure a command as follows
// <name> [OPTIONS] <Required_values> <Optional_values>
//
// the structure of this system is done to remove as much ambiuty as possible
// So as options have a - idenfying them there is no issue however due to the
// fact that optional elements can only be determined by the number of arguments
// provided (excluding flags) they are required to go at the end as otherwise
// there is no way of knowing which is which. it is also worth knowing that due
// to how the internals of our system works if there is a choice between given
// the value the commandline argument to a flag which is defined to accept values
// of a optional value the flag will allways win. this may cause confusion among
// users to please ensure that you understand how your interface will work using
// this system
func (opt *Options) Parse(args []string) (map[string]string, map[string]string, map[string]string) {
	args = args[1:]
	var (
		flags            map[string]string = make(map[string]string)
		required         map[string]string = make(map[string]string)
		optional         map[string]string = make(map[string]string)
		parsing_flags    bool              = true
		parsing_required bool              = false
		parsing_optional bool              = false
		value_index      int               = 0
	)
	h_name, _ := opt.Get_sflag('h')
	if h_name == "" {
		opt = opt.Flag("help", "Prints this help message.", false)
	} else if h_name != "help" {
		opt = opt.Long_flag("help", "Prints this help message.", false)
	}
	for index := 0; index < len(args); index++ {
		if parsing_flags {
			if len(args[index]) < 2 {
				opt.Err(fmt.Sprintf("Invalid argument %s all flags must be in the fromat -<name> or --<name>.", args[index]), "")
			}
			if args[index][0] != '-' {
				parsing_flags = false
				parsing_required = true
				index = index - 1
				if len(opt.required_values) > len(args)-index {
					opt.Err("Invalid command not all required values are given.", "Please run --help to see the structure of the command you need to run")
				}
				continue
			}
			if args[index][1] == '-' {
				if len(args[index]) < 3 {
					opt.Err(fmt.Sprintf("Invalid argument %s all flags must be in the fromat -<name> or --<name>.", args[index]), "")
				}
				name := args[index][2:]
				if opt.flags[name] == (flag_info{}) || !opt.flags[name].long {
					opt.Err(fmt.Sprintf("Invalid argument %s.", args[index]), "Please run --help to see the list of options that")
				}
				if len(args)-len(opt.required_values) > index+1 && args[index+1][0] != '-' && opt.flags[name].takes_value {
					index = index + 1
					flags[name] = args[index]
				} else {
					flags[name] = "true"
				}
			} else {
				name := args[index][1:]
				if len(name) == 1 {
					full_name, _ := opt.Get_sflag(rune(name[0]))
					if full_name == "" {
						opt.Err(fmt.Sprintf("Invalid argument %s.", args[index]), "Please run --help to see the list of options that")
					}
					if len(args)-len(opt.required_values) > index+1 && args[index+1][0] != '-' && opt.flags[full_name].takes_value {
						index = index + 1
						flags[full_name] = args[index]
					} else {
						flags[full_name] = "true"
					}
				} else {
					for _, value := range name {
						full_name, _ := opt.Get_sflag(value)
						if full_name == "" {
							opt.Err(fmt.Sprintf("Invalid argument %s.", args[index]), "Please run --help to see the list of options that")
						}
						flags[full_name] = "true"
					}
				}
			}
		} else if parsing_required {
			if value_index >= len(opt.required_values) {
				parsing_required = false
				parsing_optional = true
				index = index - 1
				continue
			}
			key := opt.required_values[value_index]
			required[key.name] = args[index]
			value_index = value_index + 1
		} else if parsing_optional {
			value_index = 0
			key := opt.optional_values[value_index]
			optional[key.name] = args[index]
			value_index = value_index + 1
		}
	}
	if flags["help"] != "" {
		opt.Help()
	}
	return flags, required, optional
}

type _parse_info struct {
	index       int
	name        string
	description string
	typ         int
}

/*
The Parsing API

This api allows for a struct to be annotated with struct tags and then used by
this method to accept the commandline arguments and parse them into a usabel
form.

The struct can therefore be annotated in the following ways
(syntax described using ebnf where key and value are string with valid values are described below)

"`gopts:", "\"", key, ["=" ,value ,] {key, [ "=", value, ] } "\"", "`"

The valid keys as given here:

	flag		This can take a value of the following "short","long" these tags describe
				the type flag it is with non value or any invalid value being understood
				being both a long and short value flag. too see the defintion for long
				and short flags read the method infomation for Options.Flag,
				Options.Short_flag and Options.Long_flag

	sflag 		takes no value and describes a short flag. see Options.Short_Flag

	lflag		takes no value and describes a long flag. see Options.Long_flag

	optional 	takes no value and describes a optional argument see Options.Optional

	required	takes no value and describes a required argument

	name		takes a single word and describes the name of the flag/parameter given

	desc | 		all text between it and the next , or the end of the tag will be
	description	taken as the description for that flag or argument used in the help text


Notes:

it should be noted that the values that are inside of the struct before being
by this method will not be changed unless a commandline argument is present and
therefore are considered default arguments.

there is also some notes on what is and is not allowed by this method:

all elements of the struct that you want to be parsed should be one of the following:
bool, Int, Int8, Int16, Int32, Int64, Uint, Uint8, Uint16, Uint32, Uint64, string
and should be exported from the struct as if they are not we will ignore them.



*/
func Parse(obj interface{}, args []string, description ...string) {
	ref := reflect.ValueOf(obj)
	val := ref.Elem()
	var desc string
	for _, val := range description {
		desc = desc + val
	}
	opts := Opts(desc)
	var info []_parse_info
	//we are looking for a series of values
	if val.Kind() != reflect.Struct {
		panic(fmt.Errorf("invalid parse type given please ensure we are parsing a struct as that is all that is supported"))
	}
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		if !field.IsExported() {
			continue
		}
		tag := field.Tag.Get("gopts")
		if tag == "" {
			continue
		}
		tags := strings.Split(tag, ",")
		var added bool = false
		for _, tag_val := range tags {
			ltag_val := strings.ToLower(tag_val)
			if strings.Contains(ltag_val, "=") {
				key_val := strings.Split(ltag_val, "=")
				if len(key_val) != 2 {
					continue
				}
				switch key_val[0] {
				case "desc", "description":
					if added {
						info[len(info)-1].description = key_val[1]
					} else {
						info = append(info, _parse_info{i, field.Name, key_val[1], required_element})
						added = true
					}
				case "name":
					if added {
						info[len(info)-1].name = key_val[1]
					} else {
						info = append(info, _parse_info{i, key_val[1], "", required_element})
						added = true
					}
				case "flag":
					switch key_val[1] {
					case "short":
						if added {
							info[len(info)-1].typ = sflag_element
						} else {
							info = append(info, _parse_info{i, field.Name, "", sflag_element})
							added = true
						}
					case "long":
						if added {
							info[len(info)-1].typ = lflag_element
						} else {
							info = append(info, _parse_info{i, field.Name, "", lflag_element})
							added = true
						}
					default:
						if added {
							info[len(info)-1].typ = flag_element
						} else {
							info = append(info, _parse_info{i, field.Name, "", flag_element})
							added = true
						}
					}
				}
			} else {
				switch ltag_val {
				case "flag":
					if added {
						info[len(info)-1].typ = flag_element
					} else {
						info = append(info, _parse_info{i, field.Name, "", flag_element})
						added = true
					}
				case "sflag", "short_flag":
					if added {
						info[len(info)-1].typ = sflag_element
					} else {
						info = append(info, _parse_info{i, field.Name, "", sflag_element})
						added = true
					}
				case "lflag", "long_flag":
					if added {
						info[len(info)-1].typ = lflag_element
					} else {
						info = append(info, _parse_info{i, field.Name, "", lflag_element})
						added = true
					}
				case "optional":
					if added {
						info[len(info)-1].typ = optional_element
					} else {
						info = append(info, _parse_info{i, field.Name, "", optional_element})
						added = true
					}
				case "required":
					if added {
						info[len(info)-1].typ = required_element
					} else {
						info = append(info, _parse_info{i, field.Name, "", required_element})
						added = true
					}
				}
			}
		}

		switch info[len(info)-1].typ {
		case required_element:
			opts = opts.Required(info[len(info)-1].name, info[len(info)-1].description)
		case optional_element:
			opts = opts.Optional(info[len(info)-1].name, info[len(info)-1].description)
		case flag_element:
			opts = opts.Flag(info[len(info)-1].name, info[len(info)-1].description, (field.Type.Kind() != reflect.Bool))
		case sflag_element:
			opts = opts.Short_flag(info[len(info)-1].name, info[len(info)-1].description, (field.Type.Kind() != reflect.Bool))
		case lflag_element:
			opts = opts.Long_flag(info[len(info)-1].name, info[len(info)-1].description, (field.Type.Kind() != reflect.Bool))
		}
	}

	flags, required, optional := opts.Parse(args)

	for _, inf := range info {
		var element_value string
		switch inf.typ {
		case required_element:
			element_value = required[inf.name]
		case optional_element:
			element_value = optional[inf.name]
		case flag_element, sflag_element, lflag_element:
			element_value = flags[inf.name]
		}
		if element_value != "" {
			switch val.Type().Field(inf.index).Type.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				ok, err := strconv.Atoi(element_value)
				if err != nil {
					opts.Err(fmt.Sprintf("Invalid default value for option %s.", strings.ToUpper(inf.name)), "Please ensure that it is an positive integer")
				}
				val.Field(inf.index).SetInt(int64(ok))
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint64:
				ok, err := strconv.Atoi(element_value)
				if err != nil {
					opts.Err(fmt.Sprintf("Invalid default value for option %s.", strings.ToUpper(inf.name)), "Please ensure that it is an integer")
				}
				val.Field(inf.index).SetUint(uint64(ok))
			case reflect.Float64, reflect.Float32:
				ok, err := strconv.ParseFloat(element_value, 64)
				if err != nil {
					opts.Err(fmt.Sprintf("Invalid default value for option %s.", strings.ToUpper(inf.name)), "Please ensure that it is an integer")
				}
				val.Field(inf.index).SetFloat(ok)
			case reflect.Bool:
				ok, err := strconv.ParseBool(element_value)
				if err != nil {
					ok = false
				}
				val.Field(inf.index).SetBool(ok)
			case reflect.String:
				val.Field(inf.index).SetString(element_value)
			}
		}
	}

}
