# (Branch: Parser) Todos

1. Type `...`
1. PatternTyped `PatternEnumVariantFn | PatternStructField | PatternBinaryTyped`
1. Pattern `PatternEnumVariantFn | PatternStructField | PatternBinary`
1. PatternEnumVariantFn `Ident LBrace PatternBinaryItem* RBrace`
1. PatternStructField `Ident LBra PatternBinaryItem* RBra`
1. PatternBinaryTyped `Ident (Colon <Type>)?`
1. PatternBinary `Ident`
1. PatternBinaryItem `Ident ','?`
1. IfElse `KIf <Expr> <Block> ( KElse (<Block> | <IfElse>) )?`
1. While `KWhile <Expr> <Block>`
1. ForIn `KFor <Pattern> KIn <Expr> <Block>`
1. Loop `Loop <Block>`
1. ExprAssign & ExprMember & ExprLookup & ExprFnCall `...`
1. Let(Pattern) `KLet <Pattern> OAssign <Expr>`
1. StructDecl `KStruct Ident LBra <StructField>* RBra`
1. StructField `Ident Colon <Type> ','?`
1. EnumDecl `KEnum Ident LBra <EnumVariant>* RBra`
1. EnumVariant `EnumVariantBinary | EnumVariantFn`
1. EnumVariantBinary `Ident ','?`
1. EnumVariantFn `Ident LBrace (EnumVariantFnArg)* RBrace ','?`
1. EnumVariantFnArg `<Type> ','?`