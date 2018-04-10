using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Threading;
using System.Diagnostics;
using Nancy.Hosting.Self;

namespace aiis
{
    class Aiis
    {
        static void Main(string[] args)
        {
            using (var nancyHost = new NancyHost(new Uri("http://0.0.0.0:8888")))
            {
                nancyHost.Start();

                try
                {
                    Process.Start("http://0.0.0.0:8888");
                }
                catch (Exception)
                {
                }
                Console.ReadKey();
            }
        }
    }
}
