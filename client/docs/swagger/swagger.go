package swagger

import (
	_ "github.com/MANTRA-Chain/mantrachain/v4/client/docs/statik" // Import MANTRA Chain statik
	"github.com/rakyll/statik/fs"
)

// https://github.com/rakyll/statik/issues/56

// FS is the MANTRA Chain swagger filesystem
var FS, _ = fs.New()
