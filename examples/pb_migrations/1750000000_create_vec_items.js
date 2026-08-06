migrate((app) => {
  app.db().newQuery(`
    CREATE VIRTUAL TABLE vec_items USING vec0(
      embedding float[3],
      +record_id TEXT
    )
  `).execute()
}, (app) => {
  app.db().newQuery("DROP TABLE IF EXISTS vec_items").execute()
})
