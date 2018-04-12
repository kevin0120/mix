using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Threading;
using System.Diagnostics;
using Nancy.Hosting.Self;
using core;

namespace aiis
{
    class Aiis
    {
        static void Main(string[] args)
        {
            // 初始化配置文件
            Console.WriteLine("init config ...");
            AiisConfig config = Config.LoadConfig();
            if(config == null)
            {
                Console.WriteLine("init config FAILED");
                return;
            }
            Console.WriteLine("init config OK");

            // 初始化日志
            Console.WriteLine("init log");

            // 初始化插件
            Console.WriteLine("init plugins ...");
            PluginSystem ps = new PluginSystem();
            ps.LoadFromConfig(config);

            HostConfiguration hc = new HostConfiguration();
            hc.UrlReservations.CreateAutomatically = true;

            string url = string.Format("http://localhost:{0}/", config.aiis.port.ToString());
            using (var nancyHost = new NancyHost(hc, new Uri(url + ApiConfig.API_TAG())))
            {
                try
                {
                    nancyHost.Start();
                    Console.WriteLine(string.Format("AIIS is running on {0}", url));
                }
                catch (Exception e)
                {
                    Console.WriteLine("AIIS FAILED:" + e.Message);
                }

                Console.ReadKey();
            }
        }

        
    }
}
