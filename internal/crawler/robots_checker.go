package crawler

import (
	"errors"
	"net/url"
	"sync"

	"github.com/stjudewashere/seonaut/internal/http_crawler"

	"github.com/temoto/robotstxt"
)

type RobotsChecker struct {
	robotsMap map[string]*robotstxt.RobotsData
	rlock     *sync.RWMutex
	client    *http_crawler.Client
}

func NewRobotsChecker(client *http_crawler.Client) *RobotsChecker {
	return &RobotsChecker{
		robotsMap: make(map[string]*robotstxt.RobotsData),
		rlock:     &sync.RWMutex{},
		client:    client,
	}
}

// Returns true if the URL is blocked by robots.txt
func (r *RobotsChecker) IsBlocked(u *url.URL) bool {
	robot, err := r.getRobotsMap(u)
	if err != nil || robot == nil {
		return false
	}

	path := u.EscapedPath()
	if u.RawQuery != "" {
		path += "?" + u.Query().Encode()
	}

	return !robot.TestAgent(path, r.client.Options.UserAgent)
}

// Returns true if the robots.txt file exists and is valid
func (r *RobotsChecker) Exists(u *url.URL) bool {
	robot, err := r.getRobotsMap(u)
	if err != nil {
		return false
	}

	if robot == nil {
		return false
	}

	return true
}

// Returns a list of sitemaps found in the robots.txt file
func (r *RobotsChecker) GetSitemaps(u *url.URL) []string {
	robot, err := r.getRobotsMap(u)
	if err != nil {
		return []string{}
	}

	return robot.Sitemaps
}

// Returns a RobotsData checking if it has already been created and stored in the robotsMap
func (r *RobotsChecker) getRobotsMap(u *url.URL) (*robotstxt.RobotsData, error) {
	r.rlock.RLock()
	robot, ok := r.robotsMap[u.Host]
	r.rlock.RUnlock()

	if !ok {
		resp, err := r.client.Get(u.Scheme + "://" + u.Host + "/robots.txt")
		if err != nil {
			r.rlock.Lock()
			r.robotsMap[u.Host] = nil
			r.rlock.Unlock()

			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			r.rlock.Lock()
			r.robotsMap[u.Host] = nil
			r.rlock.Unlock()

			return nil, errors.New("Robots.txt file does not exist")
		}

		robot, err = robotstxt.FromResponse(resp)
		if err != nil {
			r.rlock.Lock()
			r.robotsMap[u.Host] = nil
			r.rlock.Unlock()

			return nil, err
		}

		r.rlock.Lock()
		r.robotsMap[u.Host] = robot
		r.rlock.Unlock()
	}

	return robot, nil
}
