import rss from '@astrojs/rss';
import { getCollection } from 'astro:content';
export async function GET(context){const posts=(await getCollection('blog',({data})=>!data.draft)).sort((a,b)=>b.data.published.valueOf()-a.data.published.valueOf());return rss({title:"nabunana's Blog",description:'工程、Agent、底层原理与生活记录。',site:context.site,items:posts.map(post=>({title:post.data.title,description:post.data.description,pubDate:post.data.published,link:`/blog/${post.id}/`,categories:post.data.tags,author:'nabunana'}))});}
