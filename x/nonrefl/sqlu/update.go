package sqlu

//TableUpdateContext update a field based on name
/*func TableUpdateContext(ctx context.Context, db SQLer, s Schemer) (sql.Result, error) {
	schema := s.Schema()

	keys := schema.Keys()
	fields := []string{}
	params := []interface{}{}
	for _, f := range schema.Fields {
		if schema.FieldOpts != nil {
			if opt, ok := schema.FieldOpts[f.Name]; ok && opt.IsKey {
				continue
			}
		}
		fields = append(fields, f.Name+"= ?")
		params = append(params, f.Ptr)
	}

	qry := fmt.Sprintf(
		"UPDATE \"%s\" SET %s WHERE %s",
		schema.Table,
		strings.Join(fields, ", "),
		strings.Join(keys, " AND "),
	)
	log.Println("Qry:", qry, params)
	return db.ExecContext(ctx, qry, params...)
}*/
