%{
package main

import (
	  "go/ast"
	  "github.com/ds-horizon/datagen/codegen"
)
%}

/* ------------ Semantic value union ------------ */
%union{
    parsed          *codegen.DatagenParsed
    modelName       string
    fields          *ast.FieldList
    misc            string
    tags            map[string]string
    genFuns         []*codegen.GenFn
    serialiserFunc  *codegen.SerialiserFunc
    calls           []*ast.CallExpr
    count           int
    str             string
    metadata        *codegen.Metadata
}

/* ------------ Terminals (tokens) ------------ */

/* Keywords */
%token MODEL FIELDS MISC METADATA GEN_FNS CALLS FN COUNT TAGS GEN_FNS SERIALISER_FUNC

/* Punctuators */
%token L_BRACE R_BRACE L_PARENTHESIS R_PARENTHESIS COLON

/* Literals / lexeme-carrying terminals */
%token<count> COUNT_INT
%token<str>   MODEL_NAME FN_NAME FN_ARGS FN_BODY FIELDS_BODY MISC_BODY TAGS_BODY CALLS_BODY

/* ------------ Nonterminals (typed) ------------ */

%type<modelName> model_name
%type<parsed>    body

%type<fields>    fields_section
%type<misc>      misc_section
%type<metadata>  metadata_section metadata_body
%type<tags>      tags_entry
%type<genFuns>   gen_fns_section gen_fns
%type<serialiserFunc> serialiser_func_section
%type<calls>     calls_section
%type<count>     count_entry

%type<str>       fields_body misc_body tags_body calls_body


%start main

%%

main: MODEL model_name L_BRACE body R_BRACE

model_name: MODEL_NAME
{ $$ = $1 }

body: fields_section body
        {
          yylex.(*lex).parsed.Fields = $1
	    $$ = yylex.(*lex).parsed
        }
      | misc_section body
          {
	     yylex.(*lex).parsed.Misc = $1
	       $$ = yylex.(*lex).parsed
	  }
      | metadata_section body
          {
	    yylex.(*lex).parsed.Metadata = $1
	       $$ = yylex.(*lex).parsed
	  }
      | calls_section body
          {
	     yylex.(*lex).parsed.Calls = $1
	       $$ = yylex.(*lex).parsed
	  }
      | gen_fns_section body
          {
	    yylex.(*lex).parsed.GenFuns = $1
	      $$ = yylex.(*lex).parsed
	  }
      | serialiser_func_section body
          {
	    yylex.(*lex).parsed.SerialiserFunc = $1
	    $$ = yylex.(*lex).parsed
	  }
      | // empty
         {
	   $$ = yylex.(*lex).parsed
	 }

fields_section: FIELDS L_BRACE fields_body R_BRACE
{
  $$ = yylex.(*lex).parse_fields($3)
}

/* separate non-terminal to further lex the fields_body if required */
fields_body: FIELDS_BODY
{ $$ = $1 }

// misc
misc_section: MISC L_BRACE misc_body R_BRACE
{
  $$ = yylex.(*lex).parse_misc($3)
}

misc_body: MISC_BODY
{ $$ = $1 }

// metadata
metadata_section: METADATA L_BRACE metadata_body R_BRACE
{ $$ = $3 }

metadata_body: count_entry metadata_body
                {
		  if $2 == nil {
		      $2 = &codegen.Metadata{}
		  }
		  $2.Count = $1
		  $$ = $2
                }
               | tags_entry metadata_body
 	        {
		    if $2 == nil {
		        $2 = &codegen.Metadata{}
		    }
		    $2.Tags = $1
		    $$ = $2
 	        }
               | // empty
	       {}

count_entry: COUNT COLON COUNT_INT
{
  $$ = $3
}

// tags
tags_entry: TAGS COLON L_BRACE tags_body R_BRACE
{
  $$ = yylex.(*lex).parse_tags($4)
}

tags_body: TAGS_BODY
{ $$ = $1 }

// calls
calls_section: CALLS L_BRACE calls_body R_BRACE
{
  $$ = yylex.(*lex).parse_calls($3)
}

calls_body: CALLS_BODY
{ $$ = $1 }

// gens
gen_fns_section: GEN_FNS L_BRACE gen_fns R_BRACE
{
  $$ = $3
}

gen_fns: FN FN_NAME L_PARENTHESIS FN_ARGS R_PARENTHESIS L_BRACE FN_BODY R_BRACE gen_fns
            {
	        yylex.(*lex).add_gen_fn($2, $4, $7)
		$$ = yylex.(*lex).parsed.GenFuns
	    }
         | // empty
            {
              $$ = yylex.(*lex).parsed.GenFuns
            }

// serialisers
serialiser_func_section: SERIALISER_FUNC L_BRACE FN_BODY R_BRACE
{
  $$ = yylex.(*lex).add_serialiser_fn($3)
}
