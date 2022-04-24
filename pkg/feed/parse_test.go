package feed

import (
	"testing"

	"github.com/lithammer/dedent"
	"github.com/stretchr/testify/require"
)

func TestConvertStringToFeeds(t *testing.T) {
	in := dedent.Dedent(`
		http://feed.one/something/rss.xml
		http://feed.two/again/atom.xml
		http://feed.three/whatever/hopeitsxml
	`)
	feeds := convertContentsToFeeds(in)

	require.Len(t, feeds, 3)
	require.Equal(t, "http://feed.one/something/rss.xml", feeds[0].URL)
	require.Equal(t, "http://feed.two/again/atom.xml", feeds[1].URL)
	require.Equal(t, "http://feed.three/whatever/hopeitsxml", feeds[2].URL)
}
