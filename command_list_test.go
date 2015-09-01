package main

func ExampleListCommand_pretty() {
	app := makeTestEnv()
	app.Run([]string{"pgbrew", "list", "-f", "pretty"})
	// Output:
	// +-------------+----------+
	// |    NAME     | VERSION  |
	// +-------------+----------+
	// | 9.3.9-debug | REL9_3_9 |
	// | 9.4.4       | REL9_4_4 |
	// +-------------+----------+
}

func ExampleListCommand_prettyDetail() {
	app := makeTestEnv()
	app.Run([]string{"pgbrew", "list", "-f", "pretty", "-d"})
	// Output:
	// +-------------+----------+------------------------------------------+------------------------------------+---------------------------------------------+
	// |    NAME     | VERSION  |                   HASH                   |                PATH                |              CONFIGURE OPTIONS              |
	// +-------------+----------+------------------------------------------+------------------------------------+---------------------------------------------+
	// | 9.3.9-debug | REL9_3_9 | 553e576e05b50f9faffbd3dd721e44fc3746898d | /root/.pgbrew/versions/9.3.9-debug | --prefix=/root/.pgbrew/versions/9.3.9-debug |
	// |             |          |                                          |                                    | --enable-debug --enable-cassert             |
	// | 9.4.4       | REL9_4_4 | 7c055f3ec3bd338a1ebb8c73cff3d01df626471e | /root/.pgbrew/versions/9.4.4       | --prefix=/root/.pgbrew/versions/9.4.4       |
	// +-------------+----------+------------------------------------------+------------------------------------+---------------------------------------------+
}

func ExampleListCommand_plain() {
	app := makeTestEnv()
	app.Run([]string{"pgbrew", "list", "-f", "plain"})
	// Output:
	// 9.3.9-debug	REL9_3_9
	// 9.4.4	REL9_4_4
}

func ExampleListCommand_plainDetail() {
	app := makeTestEnv()
	app.Run([]string{"pgbrew", "list", "-f", "plain", "-d"})
	// Output:
	// 9.3.9-debug	REL9_3_9	553e576e05b50f9faffbd3dd721e44fc3746898d	/root/.pgbrew/versions/9.3.9-debug	--prefix=/root/.pgbrew/versions/9.3.9-debug --enable-debug --enable-cassert
	// 9.4.4	REL9_4_4	7c055f3ec3bd338a1ebb8c73cff3d01df626471e	/root/.pgbrew/versions/9.4.4	--prefix=/root/.pgbrew/versions/9.4.4
}
func ExampleListCommand_json() {
	app := makeTestEnv()
	app.Run([]string{"pgbrew", "list", "-f", "json"})
	// Output: [{"name":"9.3.9-debug","version":"REL9_3_9"},{"name":"9.4.4","version":"REL9_4_4"}]
}

func ExampleListCommand_jsonDetail() {
	app := makeTestEnv()
	app.Run([]string{"pgbrew", "list", "-f", "json", "-d"})
	// Output: [{"name":"9.3.9-debug","version":"REL9_3_9","hash":"553e576e05b50f9faffbd3dd721e44fc3746898d","path":"/root/.pgbrew/versions/9.3.9-debug","configureOptions":["--prefix=/root/.pgbrew/versions/9.3.9-debug","--enable-debug","--enable-cassert"]},{"name":"9.4.4","version":"REL9_4_4","hash":"7c055f3ec3bd338a1ebb8c73cff3d01df626471e","path":"/root/.pgbrew/versions/9.4.4","configureOptions":["--prefix=/root/.pgbrew/versions/9.4.4"]}]
}
