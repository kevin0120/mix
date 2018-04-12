using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using core;

namespace plugin
{
    public class Plugin : IPlugin
    {
        public void Init()
        {
            _fis_client = new FisClinet();
            _fis_client.Init();

            Console.WriteLine("svw init");
        }

        public void Release()
        {
            Console.WriteLine("svw release");
        }

        public void Run()
        {
            this.Init();
            Console.WriteLine("svw run");
        }

        public string Name()
        {
            return "svw";
        }

        public string Desc()
        {
            return "对接大众fis系统";
        }

        private FisClinet _fis_client;
    }
}
