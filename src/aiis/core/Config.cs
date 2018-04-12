using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using YamlDotNet.Serialization;
using YamlDotNet.Serialization.NamingConventions;

namespace core
{
    public class aiis
    {
        public uint port { get; set; }
    }

    public class log
    {
        public uint keep { get; set; }
    }

    public class plugin
    {
        public string path { get; set; }
        public bool enable { get; set; }
    }

    public class AiisConfig
    {
        public aiis aiis { get; set; }
        public log log { get; set; }
        public List<plugin> plugins { get; set; }
    }

    public class Config
    {
        public static AiisConfig LoadConfig()
        {
            AiisConfig rt = null;

            try
            {
                var file = new StringReader(System.IO.File.ReadAllText(conf_file));
                var deserializer = new DeserializerBuilder()
                    .WithNamingConvention(new CamelCaseNamingConvention())
                    .Build();
                rt = deserializer.Deserialize<AiisConfig>(file);
            }
            catch(Exception e)
            {
                Console.WriteLine("init config file error:" + e.Message);
            }

            return rt;
        }
        private static string conf_file = "./conf/aiis.conf";
    }
}
