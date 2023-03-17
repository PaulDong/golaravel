package main

func doMigrate(arg2, arg3 string) error {
	dsn := getDSN()
	// run the migration command
	switch arg2 {
	case "up":
		err := gol.MigrateUp(dsn)
		if err != nil {
			return err
		}
	case "down":
		if arg3 == "all" {
			err := gol.MigrateDownAll(dsn)
			if err != nil {
				return err
			}
		} else {
			err := gol.Steps(-1, dsn)
			if err != nil {
				return err
			}
		}
	case "reset":
		err := gol.MigrateDownAll(dsn)
		if err != nil {
			return err
		}
		err = gol.MigrateUp(dsn)
		if err != nil {
			return err
		}
	default:
		showHelp()
	}
	return nil
}