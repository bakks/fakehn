require "mechanize"
require "nokogiri"

def scrape file, url, pages = 1
  i = 0
  x = Mechanize.new

  while i < pages
    puts url

    page = x.get(url)
    titles = page.parser.xpath('//td[@class="title"]/a')

    titles[0..-2].each { |t| file.puts t.text }

    link = titles[-1]['href']
    link = '/' + link unless link[0] == '/'
    url = 'http://news.ycombinator.com' + link
    
    i += 1
  end
 
end

filename = 'titles/' + Time.now.utc.strftime('%Y%m%d.%H%M%S.txt')
file = File.open(filename, 'w')

scrape file, 'http://news.ycombinator.com', 5
scrape file, 'http://news.ycombinator.com/newest', 2

file.close
